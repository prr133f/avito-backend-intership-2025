package pr

import (
	"context"

	"github.com/prr133f/avito-backend-intership-2025/internal/domain/pr"
)

func (u usecase) Create(ctx context.Context, prModel pr.PullRequest) (pr.PullRequest, error) {
	return u.pg.PRRepo.Create(ctx, prModel)
}
func (u usecase) Merge(ctx context.Context, prId string) (pr.PullRequest, error) {
	return u.pg.PRRepo.Merge(ctx, prId)
}
func (u usecase) Reassign(ctx context.Context, prId string, oldReviewerId string) (pr.PullRequest, string, error) {
	return u.pg.PRRepo.Reassign(ctx, prId, oldReviewerId)
}
