package errors

import "github.com/pkg/errors"

var (
	ErrTeamAlreadyExists = errors.New("team already exists")
	ErrNotFound          = errors.New("resource not found")
	ErrPrAlreadyExists   = errors.New("PR is already exists")
	ErrPrMerged          = errors.New("cannot reassign on merged PR")
	ErrNotAssigned       = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate       = errors.New("no active replacement candidate in team")
)
