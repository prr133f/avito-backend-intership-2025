package user

import (
	"context"

	"github.com/prr133f/avito-backend-intership-2025/internal/domain/pr"
	"github.com/prr133f/avito-backend-intership-2025/internal/domain/user"
)

func (u usecase) SetActive(ctx context.Context, userID string, active bool) (user.User, error) {
	if err := u.pg.UserRepo.SetActive(ctx, userID, active); err != nil {
		u.log.Error("error while setting active", "err", err)
		return user.User{}, err
	}
	return u.pg.UserRepo.GetUser(ctx, userID)
}

func (u usecase) GetReview(ctx context.Context, userID string) ([]pr.PullRequest, error) {
	return u.pg.UserRepo.GetReview(ctx, userID)
}
