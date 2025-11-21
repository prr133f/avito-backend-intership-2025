package pr

import (
	"context"
	"log/slog"
	"os"

	"github.com/prr133f/avito-backend-intership-2025/internal/domain/pr"
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
	Create(ctx context.Context, pr pr.PullRequest) (pr.PullRequest, error)
	Merge(ctx context.Context, prId string) (pr.PullRequest, error)
	Reassign(ctx context.Context, prId string, oldReviewerId string) (pr.PullRequest, string, error)
}
