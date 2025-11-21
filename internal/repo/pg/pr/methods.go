package pr

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prr133f/avito-backend-intership-2025/internal/domain/pr"
	pgerrs "github.com/prr133f/avito-backend-intership-2025/internal/repo/pg/errors"
)

func (s service) Create(ctx context.Context, prModel pr.PullRequest) (pr.PullRequest, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error("error while creating pr", "err", err)
		return pr.PullRequest{}, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			s.log.Error("error while rolling back transaction", "err", err)
		}
	}()

	_, err = tx.Exec(ctx, `
	INSERT INTO pull_requests(id, name, author)
	VALUES ($1, $2, $3)`, prModel.Id, prModel.Name, prModel.Author)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.Code == "23505" {
				return pr.PullRequest{}, pgerrs.ErrPrAlreadyExists
			}
			if pgerr.Code == "23503" {
				return pr.PullRequest{}, pgerrs.ErrNotFound
			}
		}
		s.log.Error("error while creating pr", "err", err)
		return pr.PullRequest{}, err
	}

	reviewers := make([]string, 0, 2)
	rows, err := tx.Query(ctx, `
	with user_team as (
		select team_name
		from teams_users t
		where user_id = 'u1'
	)
	select tu.user_id
	from teams_users tu
	join user_team ut on tu.team_name = ut.team_name
	join users u on u.id = tu.user_id
	where
		tu.user_id != 'u1'
		and u.is_active is true`, prModel.Author)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Info("no reviewers found")
			return pr.PullRequest{}, pgerrs.ErrNotFound
		}
		s.log.Error("error while getting reviewers", "err", err)
		return pr.PullRequest{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var reviewerId string
		err = rows.Scan(&reviewerId)
		if err != nil {
			s.log.Error("error while scanning reviewer", "err", err)
			return pr.PullRequest{}, err
		}
		reviewers = append(reviewers, reviewerId)
	}

	for _, reviewer := range reviewers {
		_, err = tx.Exec(ctx, `
		INSERT INTO pull_requests_users (pull_request_id, user_id)
		VALUES ($1, $2)`, prModel.Id, reviewer)
		if err != nil {
			s.log.Error("error while adding reviewers to pr", "err", err)
			return pr.PullRequest{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("error while committing transaction", "err", err)
		return pr.PullRequest{}, err
	}

	return pr.PullRequest{
		Id:                prModel.Id,
		Name:              prModel.Name,
		Author:            prModel.Author,
		AssignedReviewers: reviewers,
		Status:            "OPEN",
	}, nil
}
func (s service) Merge(ctx context.Context, prId string) (pr.PullRequest, error) {

	return pr.PullRequest{}, nil
}
func (s service) Reassign(ctx context.Context, prId string, oldReviewerId string) (pr.PullRequest, string, error) {
	return pr.PullRequest{}, "", nil
}
