package team

import (
	"log/slog"
	"net/http"

	"github.com/prr133f/avito-backend-intership-2025/internal/rest/usecases/team"
)

type handler struct {
	log     *slog.Logger
	usecase team.Usecase
}

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

func NewHandler(log *slog.Logger) Handler {
	return &handler{
		log:     log,
		usecase: team.NewUsecase(log),
	}
}
