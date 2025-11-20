package user

import (
	"log/slog"
	"net/http"

	"github.com/prr133f/avito-backend-intership-2025/internal/usecases/user"
)

type handler struct {
	log     *slog.Logger
	usecase user.Usecase
}

type Handler interface {
	SetActive(w http.ResponseWriter, r *http.Request)
	GetReview(w http.ResponseWriter, r *http.Request)
}

func NewHandler(log *slog.Logger) Handler {
	return &handler{
		log:     log,
		usecase: user.NewUsecase(log),
	}
}
