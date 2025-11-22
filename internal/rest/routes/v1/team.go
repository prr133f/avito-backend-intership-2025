package v1

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/handlers/team"
)

func InitTeamRoutes(log *slog.Logger) http.Handler {
	r := chi.NewRouter()
	handler := team.NewHandler(log)

	r.Post("/add", handler.Create)
	r.Get("/get", handler.Get)

	return r
}
