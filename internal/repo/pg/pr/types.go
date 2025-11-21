package pr

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prr133f/avito-backend-intership-2025/internal/domain/pr"
)

type service struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

type Repo interface {
	Create(ctx context.Context, pr pr.PullRequest) (pr.PullRequest, error)
	Merge(ctx context.Context, prId string) (pr.PullRequest, error)
	Reassign(ctx context.Context, prId string, oldReviewerId string) (pr.PullRequest, string, error)
}

func NewService(log *slog.Logger, pool *pgxpool.Pool) Repo {
	return &service{
		log:  log,
		pool: pool,
	}
}
