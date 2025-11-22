package team

import (
	"context"

	"github.com/prr133f/avito-backend-intership-2025/internal/domain/team"
)

func (u usecase) Add(ctx context.Context, team team.Team) error {
	return u.pg.TeamRepo.Create(ctx, team)
}

func (u usecase) Get(ctx context.Context, name string) (team.Team, error) {
	return u.pg.TeamRepo.Get(ctx, name)
}
