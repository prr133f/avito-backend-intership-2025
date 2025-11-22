package team

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prr133f/avito-backend-intership-2025/internal/domain/team"
)

type service struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

type Repo interface {
	Create(ctx context.Context, team team.Team) error
	Get(ctx context.Context, name string) (team.Team, error)
}

func NewService(log *slog.Logger, pool *pgxpool.Pool) Repo {
	return &service{
		log:  log,
		pool: pool,
	}
}
