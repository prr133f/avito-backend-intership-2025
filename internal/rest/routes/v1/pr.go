package v1

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/handlers/pr"
)

func InitPRRoutes(log *slog.Logger) http.Handler {
	r := chi.NewRouter()
	handler := pr.NewHandler(log)

	r.Post("/create", handler.Create)
	r.Post("/merge", handler.Merge)
	r.Post("/reassign", handler.Reassign)

	return r
}
