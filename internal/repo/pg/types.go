package pg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/prr133f/avito-backend-intership-2025/internal/repo/pg/pr"
	"github.com/prr133f/avito-backend-intership-2025/internal/repo/pg/team"
	"github.com/prr133f/avito-backend-intership-2025/internal/repo/pg/user"
)

type Service struct {
	UserRepo user.Repo
	TeamRepo team.Repo
	PRRepo   pr.Repo
}

func New(logger *slog.Logger) *Service {
	pgOnce := sync.Once{}
	pgInst := &Service{}
	DSN := os.Getenv("POSTGRES_DSN")
	if DSN == "" {
		DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DB"))
	}
	pgOnce.Do(func() {
		pool, err := pgxpool.New(context.Background(), DSN)
		if err != nil {
			logger.Error("cannot connect to Service on "+DSN, slog.Any("error", err))
		}

		pgInst = &Service{
			UserRepo: user.NewService(logger, pool),
			TeamRepo: team.NewService(logger, pool),
			PRRepo:   pr.NewService(logger, pool),
		}
	})

	return pgInst
}
