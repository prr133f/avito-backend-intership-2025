package dto

type PullRequest struct {
	ID                string   `json:"pull_request_id"`
	Name              string   `json:"pull_request_name"`
	Author            string   `json:"author_id"`
	Status            string   `json:"status"`
	MergedAt          string   `json:"merged_at"`
	AssignedReviewers []string `json:"assigned_reviewers"`
}
