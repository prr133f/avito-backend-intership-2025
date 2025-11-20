package user

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prr133f/avito-backend-intership-2025/internal/domain/pr"
	"github.com/prr133f/avito-backend-intership-2025/internal/domain/user"
)

type service struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

type Repo interface {
	SetActive(ctx context.Context, userID string, active bool) error
	GetReview(ctx context.Context, userID string) ([]pr.PullRequest, error)
	GetUser(ctx context.Context, userID string) (user.User, error)
}

func NewService(log *slog.Logger, pool *pgxpool.Pool) Repo {
	return &service{
		log:  log,
		pool: pool,
	}
}
