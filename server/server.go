package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/MartyHub/mail-api/db"
	"github.com/MartyHub/mail-api/health"
	"github.com/MartyHub/mail-api/mail"
	"github.com/MartyHub/mail-api/version"
	"github.com/invopop/validation"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	cfg Config
	e   *echo.Echo
}

func NewServer(cfg Config) Server {
	result := Server{
		cfg: cfg,
		e:   echo.New(),
	}

	result.e.Debug = cfg.Development
	result.e.DisableHTTP2 = true
	result.e.HideBanner = true
	result.e.HTTPErrorHandler = result.newErrorHandler()
	result.e.Server.ReadTimeout = cfg.ReadTimeout
	result.e.Server.WriteTimeout = cfg.WriteTimeout
	result.e.Validator = validator{}

	result.e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogError:   true,
		LogLatency: true,
		LogMethod:  true,
		LogStatus:  true,
		LogURI:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info().
				Err(v.Error).
				Str("latency", v.Latency.String()).
				Str("method", v.Method).
				Int("status", v.Status).
				Str("URI", v.URI).
				Msg("request")

			return nil
		},
	}))

	return result
}

func (s Server) Routes(repo db.Repository) {
	healthHandler := health.NewHandler(repo)
	mailHandler := mail.NewHandler(repo)
	versionHandler := version.NewHandler()

	root := s.e.Group("/api/v1")
	root.GET("/health", healthHandler.Get)
	root.GET("/version", versionHandler.Get)

	gmail := root.Group("/mail")
	gmail.POST("", mailHandler.Create)
	gmail.GET("/:id", mailHandler.Get)
}

func (s Server) Start() {
	go func() {
		if err := s.e.Start(s.cfg.ServerAddress()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		log.Info().Msgf("Received QUIT signal")
		s.cfg.Sender.Stop()
	case <-s.cfg.Sender.Stopper:
	}

	s.stop()
}

func (s Server) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	log.Info().Msgf("Shutting down server...")

	if err := s.e.Shutdown(ctx); err != nil {
		log.Err(err).Msg("Failed to shutdown server")
	}

	log.Info().Msgf("Waiting for senders...")

	s.cfg.Sender.Waiter.Wait()

	log.Info().Msgf("Exit")
}

func (s Server) newErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var ve *validation.Errors
		if errors.As(err, &ve) {
			if err = c.JSON(http.StatusBadRequest, ve); err != nil {
				c.Logger().Error(err)
			}

			return
		}

		if errors.Is(err, pgx.ErrNoRows) {
			c.Error(echo.ErrNotFound)

			return
		}

		s.e.DefaultHTTPErrorHandler(err, c)
	}
}
