package health

import (
	"context"
	"time"

	"github.com/MartyHub/mail-api/db"
)

const statusOK = "OK"

type health struct {
	Database string    `json:"database"`
	Time     time.Time `json:"time"`
}

type Service interface {
	Health(ctx context.Context) health
}

func newService(repo db.Repository) Service {
	return service{repo: repo}
}

type service struct {
	repo db.Repository
}

func (s service) Health(ctx context.Context) health {
	now := s.repo.Now()
	result := health{
		Database: statusOK,
		Time:     now.Time,
	}

	if err := s.repo.Ping(ctx); err != nil {
		result.Database = err.Error()
	}

	return result
}
