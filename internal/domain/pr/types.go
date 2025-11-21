package pr

import "time"

type PullRequest struct {
	Id                string     `db:"id"`
	Name              string     `db:"name"`
	Author            string     `db:"author"`
	Status            string     `db:"status"`
	AssignedReviewers []string   `db:"assigned_reviewers"`
	MergedAt          *time.Time `db:"merged_at"`
}
