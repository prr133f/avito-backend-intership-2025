package user

import (
	"encoding/json"
	"errors"
	"net/http"

	pgerrs "github.com/prr133f/avito-backend-intership-2025/internal/repo/pg/errors"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/dto"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/handlers"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/mapper"
)

func (h handler) SetActive(w http.ResponseWriter, r *http.Request) {
	var dtouser dto.User
	if err := json.NewDecoder(r.Body).Decode(&dtouser); err != nil {
		h.log.Error("error while parsing body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	domainUser, err := h.usecase.SetActive(r.Context(), dtouser.Id, dtouser.IsActive)
	if err != nil {
		if errors.Is(err, pgerrs.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"error": dto.Error{
					Code:    handlers.CODE_NOT_FOUND,
					Message: err.Error(),
				},
			}); err != nil {
				h.log.Error("error while encoding error response", "err", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error("error while setting active", "err", err)
		return
	}

	user := mapper.UserDomainToDTO(domainUser)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"user": user,
	}); err != nil {
		h.log.Error("error while encoding response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h handler) GetReview(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	prs, err := h.usecase.GetReview(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.log.Error("error while getting review", "err", err)
		return
	}

	mappedPrs := make([]dto.PullRequest, len(prs))
	for i, pr := range prs {
		mappedPrs[i] = mapper.PullRequestToDTO(pr)
	}
	if err := json.NewEncoder(w).Encode(map[string]any{
		"user_id":       userId,
		"pull_requests": mappedPrs,
	}); err != nil {
		h.log.Error("error while encoding response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
