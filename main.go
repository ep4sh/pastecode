package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pastecode/pkg/app"
	"pastecode/pkg/handlers"

	"go.uber.org/zap"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Minute
	idleTimeout  = 30 * time.Second
)

func main() {
	backendApp := app.NewApp()

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         backendApp.WebserverConf.Addr(),
		Handler:      handlers.LoggingMiddleware(backendApp, mux),
		ErrorLog:     zap.NewStdLog(backendApp.Sugar.Desugar()),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	// IMPORTANT!: More specific routes are higher in order
	mux.HandleFunc("/paste/{uuid}", handlers.Paste(backendApp))
	mux.HandleFunc("/add", handlers.Add(backendApp))

	mux.HandleFunc("/static/", handlers.Static(backendApp))

	// System and metrics endpoints.
	mux.HandleFunc("/healthz", handlers.Healthz(backendApp))
	mux.HandleFunc("/readyz", handlers.Readyz(backendApp))

	// Index
	mux.HandleFunc("/", handlers.Index(backendApp))
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if err != http.ErrServerClosed {
				backendApp.Sugar.Fatalf("Err occured while starting server: %v", err)
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if err := server.Shutdown(backendApp.Ctx); err != nil {
		backendApp.Sugar.Fatalf("HTTP server shutdown error: %v", err)
	}
	backendApp.Sugar.Infof("Graceful shutdown complete.")
}
