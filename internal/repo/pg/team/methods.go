package team

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prr133f/avito-backend-intership-2025/internal/domain/team"
	"github.com/prr133f/avito-backend-intership-2025/internal/domain/user"
	pgerrs "github.com/prr133f/avito-backend-intership-2025/internal/repo/pg/errors"
)

func (s service) Create(ctx context.Context, team team.Team) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error("creating tx error", "err", err)
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			s.log.Error("error rolling back tx", "err", err)
		}
	}()

	_, err = tx.Exec(ctx, `
	INSERT INTO teams (name)
	VALUES ($1)`, team.Name)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.Code == "23505" {
				return pgerrs.ErrTeamAlreadyExists
			}
		}
		s.log.Error("error creating team", "err", err)
		return err
	}

	for _, user := range team.Users {
		_, err := tx.Exec(ctx, `
		INSERT INTO users (id, name, is_active)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
			SET name = $2,
				is_active = $3`, user.Id, user.Name, user.IsActive)
		if err != nil {
			s.log.Error("error creating user in team", "err", err)
			return err
		}
		_, err = tx.Exec(ctx, `
		INSERT INTO teams_users (team_name, user_id)
		VALUES ($1, $2)
		ON CONFLICT (team_name, user_id) DO NOTHING`, team.Name, user.Id)
		if err != nil {
			s.log.Error("error creating team-user relation", "err", err)
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("error committing tx", "err", err)
		return err
	}

	return nil
}

func (s service) Get(ctx context.Context, name string) (team.Team, error) {
	var t team.Team

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Error("error starting tx", "err", err)
		return team.Team{}, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			s.log.Error("error rolling back tx", "err", err)
		}
	}()

	row := tx.QueryRow(ctx, `
	SELECT name
	FROM teams
	WHERE name = $1`, name)
	if err := row.Scan(&t.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return team.Team{}, pgerrs.ErrNotFound
		}
		s.log.Error("error scanning team", "err", err)
		return team.Team{}, err
	}

	rows, err := tx.Query(ctx, `
	SELECT u.id, u.name, u.is_active
	FROM users u
	JOIN teams_users tu ON u.id = tu.user_id
	WHERE tu.team_name = $1`, t.Name)
	if err != nil {
		s.log.Error("error while getting team users", "err", err)
		return team.Team{}, err
	}

	for rows.Next() {
		var user user.User
		if err := rows.Scan(&user.Id, &user.Name, &user.IsActive); err != nil {
			s.log.Error("error scanning user", "err", err)
			return team.Team{}, err
		}
		t.Users = append(t.Users, user)
	}
	if err := rows.Err(); err != nil {
		s.log.Error("error iterating over team users", "err", err)
		return team.Team{}, err
	}

	return t, nil
}
