package pg

import (
	"context"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prr133f/avito-backend-intership-2025/internal/repo/pg/team"
	"github.com/prr133f/avito-backend-intership-2025/internal/repo/pg/user"
)

type Service struct {
	UserRepo user.Repo
	TeamRepo team.Repo
}

func New(DSN string, logger *slog.Logger) *Service {
	pgOnce := sync.Once{}
	pgInst := &Service{}

	pgOnce.Do(func() {
		pool, err := pgxpool.New(context.Background(), DSN)
		if err != nil {
			logger.Error("cannot connect to Service on "+DSN, slog.Any("error", err))
		}

		pgInst = &Service{
			UserRepo: user.NewService(logger, pool),
			TeamRepo: team.NewService(logger, pool),
		}
	})

	return pgInst
}
