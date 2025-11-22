package user

import (
	"context"
	"log/slog"
	"os"

	"github.com/prr133f/avito-backend-intership-2025/internal/domain/pr"
	"github.com/prr133f/avito-backend-intership-2025/internal/domain/user"
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
	SetActive(ctx context.Context, userID string, active bool) (user.User, error)
	GetReview(ctx context.Context, userID string) ([]pr.PullRequest, error)
}
