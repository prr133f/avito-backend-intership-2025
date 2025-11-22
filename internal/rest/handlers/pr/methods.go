package pr

import (
	"encoding/json"
	"net/http"

	pgerrs "github.com/prr133f/avito-backend-intership-2025/internal/repo/pg/errors"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/dto"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/handlers"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/mapper"
)

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	var in dto.PullRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pr, err := h.usecase.Create(r.Context(), mapper.PullRequestFromDTO(in))
	if err != nil {
		switch err {
		case pgerrs.ErrPrAlreadyExists:
			w.WriteHeader(http.StatusConflict)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"error": dto.Error{
					Code:    handlers.CODE_PULL_REQUEST_ALREADY_EXISTS,
					Message: err.Error(),
				},
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		case pgerrs.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"error": dto.Error{
					Code:    handlers.CODE_NOT_FOUND,
					Message: err.Error(),
				},
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(mapper.PullRequestToDTO(pr)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (h handler) Merge(w http.ResponseWriter, r *http.Request) {
	var in dto.PullRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pr, err := h.usecase.Merge(r.Context(), mapper.PullRequestFromDTO(in).Id)
	if err != nil {
		switch err {
		case pgerrs.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"error": dto.Error{
					Code:    handlers.CODE_NOT_FOUND,
					Message: err.Error(),
				},
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(mapper.PullRequestToDTO(pr)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (h handler) Reassign(w http.ResponseWriter, r *http.Request) {
	var in dto.ReassignPullRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pr, newReviewerId, err := h.usecase.Reassign(r.Context(), in.ID, in.OldReviewer)
	if err != nil {
		switch err {
		case pgerrs.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"error": dto.Error{
					Code:    handlers.CODE_NOT_FOUND,
					Message: err.Error(),
				},
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		case pgerrs.ErrPrMerged:
			w.WriteHeader(http.StatusConflict)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"error": dto.Error{
					Code:    handlers.CODE_PR_MERGED,
					Message: err.Error(),
				},
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		case pgerrs.ErrNoCandidate:
			w.WriteHeader(http.StatusConflict)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"error": dto.Error{
					Code:    handlers.CODE_NO_CANDIDATE,
					Message: err.Error(),
				},
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		case pgerrs.ErrNotAssigned:
			w.WriteHeader(http.StatusConflict)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"error": dto.Error{
					Code:    handlers.CODE_NOT_ASSIGNED,
					Message: err.Error(),
				},
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"pr":          mapper.PullRequestToDTO(pr),
		"replaced_by": newReviewerId,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
