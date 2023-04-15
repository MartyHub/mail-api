package server

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/MartyHub/mail-api/db"
	"github.com/MartyHub/mail-api/smtp"
	"github.com/MartyHub/mail-api/worker"
	"github.com/caarlos0/env/v8"
	"github.com/invopop/validation"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Development bool

	Host string `envDefault:"localhost"`
	Port int    `envDefault:"8125"`

	ReadTimeout     time.Duration `envDefault:"10s"`
	ShutdownTimeout time.Duration `envDefault:"10s"`
	WriteTimeout    time.Duration `envDefault:"10s"`

	Database db.Config     `envPrefix:"DATABASE_"`
	Sender   worker.Config `envPrefix:"SENDER_"`
	SMTP     smtp.Config   `envPrefix:"SMTP_"`
}

func (c Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required, validation.Min(1)),
		validation.Field(&c.ReadTimeout, validation.Required),
		validation.Field(&c.ShutdownTimeout, validation.Required),
		validation.Field(&c.WriteTimeout, validation.Required),
		validation.Field(&c.Database),
		validation.Field(&c.Sender),
		validation.Field(&c.SMTP),
	)
}

func (c Config) String() string {
	sb := strings.Builder{}

	sb.WriteString("Config:\n")
	sb.WriteString(fmt.Sprintf("  - Development: %v\n", c.Development))
	sb.WriteString(fmt.Sprintf("  - Host: %s\n", c.Host))
	sb.WriteString(fmt.Sprintf("  - Port: %d\n", c.Port))
	sb.WriteString(fmt.Sprintf("  - Read Timeout: %v\n", c.ReadTimeout))
	sb.WriteString(fmt.Sprintf("  - Shutdown Timeout: %v\n", c.ShutdownTimeout))
	sb.WriteString(fmt.Sprintf("  - Write Timeout: %v\n", c.WriteTimeout))
	sb.WriteString(c.Sender.String())
	sb.WriteString(c.SMTP.String())

	return sb.String()
}

func ParseConfig() (Config, error) {
	var result Config

	if err := env.ParseWithOptions(&result, env.Options{
		Prefix:                "MAIL_API_",
		UseFieldNameByDefault: true,
	}); err != nil {
		return result, err
	}

	result.Sender.Stopper = make(chan bool, 1)
	result.Sender.Waiter = &sync.WaitGroup{}

	if result.Development {
		zerolog.TimeFieldFormat = time.RFC3339Nano

		cw := zerolog.ConsoleWriter{Out: os.Stderr}
		cw.TimeFormat = "15:04:05.000"

		log.Logger = log.Output(cw)
	}

	log.Info().Msg(result.String())

	return result, result.Validate()
}
