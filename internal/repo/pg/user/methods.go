package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/prr133f/avito-backend-intership-2025/internal/domain/pr"
	"github.com/prr133f/avito-backend-intership-2025/internal/domain/user"
)

func (s service) SetActive(ctx context.Context, userID string, active bool) error {
	_, err := s.pool.Exec(ctx, `
	UPDATE users
	SET is_active = $1
	WHERE id = $2`, active, userID)
	if err != nil {
		s.log.Error("setting user active err", "err", err)
		return err
	}
	return nil
}

func (s service) GetReview(ctx context.Context, userID string) ([]pr.PullRequest, error) {
	rows, err := s.pool.Query(ctx, `
	SELECT p.id, p.name, p.author, p.status
	FROM pull_requests p
	JOIN pull_requests_users pu ON p.id = pu.pull_request_id
	WHERE pu.user_id = $1`, userID)
	if err != nil {
		s.log.Error("getting user reviews err", "err", err)
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[pr.PullRequest])
}

func (s service) GetUser(ctx context.Context, userID string) (user.User, error) {
	rows, err := s.pool.Query(ctx, `
	SELECT u.id, u.name, u.is_active, tu.team_name
	FROM users u
	JOIN teams_users tu ON u.id = tu.user_id
	WHERE u.id = $1`, userID)
	if err != nil {
		s.log.Error("getting user err", "err", err)
		return user.User{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[user.User])
}
