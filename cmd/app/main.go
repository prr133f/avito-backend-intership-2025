package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prr133f/avito-backend-intership-2025/internal/logging"
	"github.com/prr133f/avito-backend-intership-2025/internal/rest/routes"
)

func main() {
	log := logging.NewLogger()
	router := routes.InitRouter(log)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("APP_PORT"),
		Handler: router,
	}

	go func() {
		log.Info("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen error", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info("shutdown signal received", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "err", err)
	}

	log.Info("server exiting")
}
