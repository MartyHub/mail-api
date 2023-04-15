package worker

import (
	"context"
	"time"

	"github.com/MartyHub/mail-api/db"
	gdb "github.com/MartyHub/mail-api/db/gen"
	"github.com/MartyHub/mail-api/smtp"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

type Sender interface {
	Run()
}

func NewSender(
	id int,
	cfg Config,
	repo db.Repository,
	smtpService smtp.Service,
) Sender {
	return sender{
		id:          id,
		cfg:         cfg,
		repo:        repo,
		smtpService: smtpService,
	}
}

type sender struct {
	id          int
	cfg         Config
	repo        db.Repository
	smtpService smtp.Service
}

func (s sender) Run() {
	s.cfg.Waiter.Add(1)

	ticker := time.NewTicker(s.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.cfg.Stopper:
			log.Info().Int("sender", s.id).Msgf("Received stop signal")
			s.cfg.Waiter.Done()

			return
		case <-ticker.C:
			s.checkMails(context.Background())
		}
	}
}

func (s sender) checkMails(ctx context.Context) {
	err := s.repo.Wrap(ctx, db.TxWrite(), func(querier gdb.Querier) error {
		log.Info().Int("sender", s.id).Msgf("Checking for mails to send...")

		mails, err := querier.GetMailsToSend(ctx, gdb.GetMailsToSendParams{
			RetryAt: s.repo.Now(),
			Limit:   s.cfg.BatchSize,
		})
		if err != nil {
			return err
		}

		log.Info().Int("sender", s.id).Msgf("Found %d mail(s) to send", len(mails))

		for _, m := range mails {
			newStatus := s.send(m)
			now := s.repo.Now()

			if err = querier.UpdateMail(ctx, gdb.UpdateMailParams{
				ID:        m.ID,
				Status:    newStatus,
				UpdatedAt: now,
				RetryAt: pgtype.Timestamp{
					Time:  s.computeRetryAt(now.Time, m.Try+1),
					Valid: true,
				},
			}); err != nil {
				log.Err(err).Int("sender", s.id).Msgf("Failed to update mail # %d to status %v", m.ID, newStatus)

				s.cfg.Stopper <- true

				break
			}
		}

		return nil
	})
	if err != nil {
		log.Err(err).Int("sender", s.id).Msg("Failed to check mails")
	}
}

func (s sender) computeRetryAt(now time.Time, try int16) time.Time {
	return now.Add(s.cfg.RetryDelay * time.Duration(try*try))
}

func (s sender) send(m gdb.Mail) gdb.MailStatus {
	try := m.Try + 1
	log.Info().Int("sender", s.id).Msgf("Sending mail # %d (try %d)...", m.ID, try)

	if err := s.smtpService.Send(smtp.Mail{
		From:    m.FromAddress,
		To:      m.RecipientsTo,
		Cc:      m.RecipientsCc,
		Subject: m.Subject,
		Body:    m.Body,
		HTML:    m.Html,
	}); err != nil {
		log.Err(err).Int("sender", s.id).Msgf("Failed to send mail # %d (try %d)", m.ID, try)

		if m.Try == s.cfg.MaxTries-1 {
			log.Info().Int("sender", s.id).Msgf("Mail # %d has reached max tries", m.ID)

			return gdb.MailStatusERROR
		}

		log.Info().Int("sender", s.id).Msgf("Mail # %d will go back in the queue", m.ID)

		return gdb.MailStatusQUEUED
	}

	log.Info().Int("sender", s.id).Msgf("Mail # %d successfully sent", m.ID)

	return gdb.MailStatusSENT
}
