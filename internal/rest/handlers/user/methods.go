package user

import (
	"encoding/json"
	"net/http"

	"github.com/prr133f/avito-backend-intership-2025/internal/rest/dto"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/mapper"
)

func (h handler) SetActive(w http.ResponseWriter, r *http.Request) {
	var dto dto.User
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.log.Error("error while parsing body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	domainUser, err := h.usecase.SetActive(r.Context(), dto.Id, dto.IsActive)
	if err != nil {
		h.log.Error("error while setting active", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func (h handler) GetReview(w http.ResponseWriter, r *http.Request) {}
