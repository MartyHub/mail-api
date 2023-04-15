package main

import (
	"fmt"

	"github.com/MartyHub/mail-api/db"
	"github.com/MartyHub/mail-api/server"
	"github.com/MartyHub/mail-api/smtp"
	"github.com/MartyHub/mail-api/utils"
	"github.com/MartyHub/mail-api/worker"
)

func main() {
	cfg, err := server.ParseConfig()
	if err != nil {
		panic(fmt.Errorf("failed to parse config: %w", err))
	}

	clock := utils.UTCClock{}

	repo, err := db.NewRepository(cfg.Database, clock)
	if err != nil {
		panic(fmt.Errorf("failed to init DB: %w", err))
	}

	defer repo.Close()

	smtpService, err := smtp.NewService(cfg.SMTP)
	if err != nil {
		panic(fmt.Errorf("failed to init SMTP: %w", err))
	}

	s := server.NewServer(cfg)
	s.Routes(repo)

	for i := 1; i <= cfg.Sender.Count; i++ {
		go worker.NewSender(i, cfg.Sender, repo, smtpService).
			Run()
	}

	s.Start()
}
