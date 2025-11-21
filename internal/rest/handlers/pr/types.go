package pr

import (
	"log/slog"
	"net/http"

	"github.com/prr133f/avito-backend-intership-2025/internal/rest/usecases/pr"
)

type handler struct {
	log     *slog.Logger
	usecase pr.Usecase
}

type Handler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Merge(w http.ResponseWriter, r *http.Request)
	Reassign(w http.ResponseWriter, r *http.Request)
}

func NewHandler(log *slog.Logger) Handler {
	return &handler{
		log:     log,
		usecase: pr.NewUsecase(log),
	}
}
