package mapper

import (
	"time"

	"github.com/prr133f/avito-backend-intership-2025/internal/domain/pr"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/dto"
)

func PullRequestToDTO(pr pr.PullRequest) dto.PullRequest {
	formatMergedAt := ""
	if pr.MergedAt != nil {
		formatMergedAt = pr.MergedAt.Format(time.RFC3339)
	}
	return dto.PullRequest{
		ID:                pr.Id,
		Name:              pr.Name,
		Author:            pr.Author,
		Status:            pr.Status,
		MergedAt:          formatMergedAt,
		AssignedReviewers: pr.AssignedReviewers,
	}
}

func PullRequestFromDTO(dto dto.PullRequest) pr.PullRequest {
	return pr.PullRequest{
		Id:                dto.ID,
		Name:              dto.Name,
		Author:            dto.Author,
		Status:            dto.Status,
		AssignedReviewers: dto.AssignedReviewers,
	}
}
