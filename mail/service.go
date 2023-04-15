package mail

import (
	"context"
	"time"

	"github.com/MartyHub/mail-api/db"
	gdb "github.com/MartyHub/mail-api/db/gen"
	"github.com/invopop/validation"
)

type Service interface {
	Create(ctx context.Context, input CreateInput) (CreateOutput, error)
	Get(ctx context.Context, input GetInput) (GetOutput, error)
}

func newService(repo db.Repository) Service {
	return &service{repo: repo}
}

type service struct {
	repo db.Repository
}

type CreateInput struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	HTML    bool     `json:"html"`
}

func (input CreateInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.From, validation.Required),
		validation.Field(&input.To, validation.Required),
		validation.Field(&input.Subject, validation.Required),
		validation.Field(&input.Body, validation.Required),
	)
}

type CreateOutput struct {
	ID        int32     `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s service) Create(ctx context.Context, input CreateInput) (CreateOutput, error) {
	var result CreateOutput

	err := s.repo.Wrap(ctx, db.TxWrite(), func(querier gdb.Querier) error {
		now := s.repo.Now()
		entity, err := querier.CreateMail(ctx, gdb.CreateMailParams{
			FromAddress:  input.From,
			RecipientsTo: input.To,
			RecipientsCc: input.Cc,
			Subject:      input.Subject,
			Body:         input.Body,
			Html:         input.HTML,
			Status:       gdb.MailStatusQUEUED,
			CreatedAt:    now,
			UpdatedAt:    now,
			RetryAt:      now,
		})
		if err != nil {
			return err
		}

		result.ID = entity.ID
		result.Status = string(entity.Status)
		result.CreatedAt = entity.CreatedAt.Time

		return nil
	})

	return result, err
}

type GetInput struct {
	ID int32 `param:"id"`
}

func (input GetInput) Validate() error {
	return validation.ValidateStruct(&input,
		validation.Field(&input.ID, validation.Required),
	)
}

type GetOutput struct {
	ID        int32     `json:"id"`
	From      string    `json:"from"`
	To        []string  `json:"to"`
	Cc        []string  `json:"cc"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	HTML      bool      `json:"html"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	RetryAt   time.Time `json:"retryAt"`
	Try       int16     `json:"try"`
}

func (s service) Get(ctx context.Context, input GetInput) (GetOutput, error) {
	var result GetOutput

	err := s.repo.Wrap(ctx, db.TxReadOnly(), func(querier gdb.Querier) error {
		entity, err := querier.GetMail(ctx, input.ID)
		if err != nil {
			return err
		}

		result.ID = entity.ID
		result.From = entity.FromAddress
		result.To = entity.RecipientsTo
		result.Cc = entity.RecipientsCc
		result.Subject = entity.Subject
		result.Body = entity.Body
		result.HTML = entity.Html
		result.Status = string(entity.Status)
		result.CreatedAt = entity.CreatedAt.Time
		result.UpdatedAt = entity.UpdatedAt.Time
		result.RetryAt = entity.RetryAt.Time
		result.Try = entity.Try

		return nil
	})

	return result, err
}
