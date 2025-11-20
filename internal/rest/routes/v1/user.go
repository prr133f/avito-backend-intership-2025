package v1

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/prr133f/avito-backend-intership-2025/internal/rest/handlers/user"
)

func InitUserRoutes(log *slog.Logger) http.Handler {
	r := chi.NewRouter()
	handler := user.NewHandler(log)

	r.Post("/setIsActive", handler.SetActive)
	r.Get("/getReview", handler.GetReview)

	return r
}
