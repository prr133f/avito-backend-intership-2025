package pr

type PullRequest struct {
	Id     string `db:"id"`
	Name   string `db:"name"`
	Author int    `db:"author"`
	Status string `db:"status"`
}
