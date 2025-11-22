package routes

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	v1 "github.com/prr133f/avito-backend-intership-2025/internal/rest/routes/v1"
)

func InitRouter(log *slog.Logger) *chi.Mux {
	r := chi.NewMux()

	r.Use(middleware.Logger)

	r.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Mount("/team", v1.InitTeamRoutes(log))

	r.Mount("/user", v1.InitUserRoutes(log))

	r.Mount("/pullRequest", v1.InitPRRoutes(log))

	return r
}
