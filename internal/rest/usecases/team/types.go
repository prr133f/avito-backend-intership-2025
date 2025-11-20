package team

import (
	"context"
	"log/slog"
	"os"

	"github.com/prr133f/avito-backend-intership-2025/internal/domain/team"
	"github.com/prr133f/avito-backend-intership-2025/internal/repo/pg"
)

type usecase struct {
	log *slog.Logger
	pg  *pg.Service
}

func NewUsecase(log *slog.Logger) Usecase {
	return &usecase{
		log: log,
		pg:  pg.New(os.Getenv("POSTGRES_DSN"), log),
	}
}

type Usecase interface {
	Add(ctx context.Context, team team.Team) error
	Get(ctx context.Context, name string) (team.Team, error)
}
