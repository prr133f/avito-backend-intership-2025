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
		where user_id = $1
	)
	select tu.user_id
	from teams_users tu
	join user_team ut on tu.team_name = ut.team_name
	join users u on u.id = tu.user_id
	where
		tu.user_id != $1
		and u.is_active is true
	limit 2`, prModel.Author)
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
	ct, err := s.pool.Exec(ctx, `
	UPDATE pull_requests
	SET status='MERGED',
		merged_at=CASE
		WHEN status != 'MERGED' THEN NOW() AT TIME ZONE 'Europe/Moscow'
		ELSE merged_at
	END
	WHERE id=$1`, prId)
	if err != nil {
		s.log.Error("error while merging pr", "err", err)
		return pr.PullRequest{}, err
	}
	if ct.RowsAffected() == 0 {
		return pr.PullRequest{}, pgerrs.ErrNotFound
	}

	return s.getPullRequest(ctx, prId)
}
func (s service) Reassign(ctx context.Context, prId string, oldReviewerId string) (pr.PullRequest, string, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error("error while starting tx for reassign", "err", err)
		return pr.PullRequest{}, "", err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			s.log.Error("error while rolling back tx in reassign", "err", err)
		}
	}()

	var status string
	if err := tx.QueryRow(ctx, `SELECT status FROM pull_requests WHERE id = $1`, prId).Scan(&status); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr.PullRequest{}, "", pgerrs.ErrNotFound
		}
		s.log.Error("error while selecting pr status", "err", err)
		return pr.PullRequest{}, "", err
	}
	if status == "MERGED" {
		return pr.PullRequest{}, "", pgerrs.ErrPrMerged
	}

	var tmp int
	if err := tx.QueryRow(ctx, `SELECT 1 FROM users WHERE id = $1`, oldReviewerId).Scan(&tmp); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr.PullRequest{}, "", pgerrs.ErrNotFound
		}
		s.log.Error("error while checking old reviewer existence", "err", err)
		return pr.PullRequest{}, "", err
	}

	if err := tx.QueryRow(ctx, `SELECT 1 FROM pull_requests_users WHERE pull_request_id = $1 AND user_id = $2`, prId, oldReviewerId).Scan(&tmp); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr.PullRequest{}, "", pgerrs.ErrNotAssigned
		}
		s.log.Error("error while checking old reviewer assigned", "err", err)
		return pr.PullRequest{}, "", err
	}

	var candidate string
	err = tx.QueryRow(ctx, `
	WITH old_reviewer_team AS (
	    SELECT team_name
	    FROM teams_users tu
	    WHERE tu.user_id = $1
	    LIMIT 1
	)
	SELECT u.id
	FROM users u
	JOIN teams_users tu ON u.id = tu.user_id
	JOIN old_reviewer_team ort ON tu.team_name = ort.team_name
	LEFT JOIN pull_requests_users pru ON pru.user_id = u.id AND pru.pull_request_id = $2
	WHERE u.is_active = TRUE
	  AND u.id != $1
	  AND pru.user_id IS NULL
	LIMIT 1`, oldReviewerId, prId).Scan(&candidate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr.PullRequest{}, "", pgerrs.ErrNoCandidate
		}
		s.log.Error("error while selecting candidate", "err", err)
		return pr.PullRequest{}, "", err
	}

	ct, err := tx.Exec(ctx, `UPDATE pull_requests_users SET user_id = $1 WHERE pull_request_id = $2 AND user_id = $3`, candidate, prId, oldReviewerId)
	if err != nil {
		s.log.Error("error while updating pull_requests_users", "err", err)
		return pr.PullRequest{}, "", err
	}
	if ct.RowsAffected() == 0 {
		return pr.PullRequest{}, "", pgerrs.ErrNotAssigned
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("error while committing reassign tx", "err", err)
		return pr.PullRequest{}, "", err
	}

	pull, err := s.getPullRequest(ctx, prId)
	if err != nil {
		return pr.PullRequest{}, "", err
	}
	return pull, candidate, nil
}

func (s service) getPullRequest(ctx context.Context, id string) (pr.PullRequest, error) {
	rows, err := s.pool.Query(ctx, `
	select
		p.id,
		p.name,
		p.author,
		p.status,
		p.merged_at,
		array_agg(pu.user_id order by pu.user_id) as assigned_reviewers
	from pull_requests p
	join pull_requests_users pu on p.id = pu.pull_request_id
	where p.id=$1
	group by p.id, p.name, p.author, p.status`, id)
	if err != nil {
		s.log.Error("error while scanning author", "err", err)
		return pr.PullRequest{}, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByNameLax[pr.PullRequest])
}
