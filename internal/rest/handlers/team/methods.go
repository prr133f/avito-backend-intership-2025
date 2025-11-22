package team

import (
	"encoding/json"
	"errors"
	"net/http"

	pgerrs "github.com/prr133f/avito-backend-intership-2025/internal/repo/pg/errors"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/dto"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/handlers"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/mapper"
)

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	var team dto.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		h.log.Error("error while decoding body", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	domainTeam, err := mapper.TeamDTOToDomain(team)
	if err != nil {
		h.log.Error("error while converting DTO to domain", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.usecase.Add(r.Context(), domainTeam); err != nil {
		if errors.Is(err, pgerrs.ErrTeamAlreadyExists) {
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string]any{
				"error": dto.Error{
					Code:    handlers.CODE_TEAM_ALREADY_EXISTS,
					Message: err.Error(),
				},
			}); err != nil {
				h.log.Error("error while encoding error response", "err", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		h.log.Error("error while creating team", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
func (h handler) Get(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")

	domainTeam, err := h.usecase.Get(r.Context(), teamName)
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
		h.log.Error("error while getting team", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	team := mapper.TeamDomainToDTO(domainTeam)

	if err := json.NewEncoder(w).Encode(team); err != nil {
		h.log.Error("error while encoding team", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
