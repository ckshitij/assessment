package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"goservice/configs"
	"goservice/internal/client"
	"goservice/internal/student"
	"log"
)

func main() {

	conf := configs.Load()

	backend := client.NewBackendClient(conf.NodeServer.BaseURL)
	studentsrv := student.NewService(backend, conf.NodeServer.Username, conf.NodeServer.Password)
	handler := student.NewHandler(studentsrv)

	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Mount("/api/v1", handler.Routes())

	addr := fmt.Sprintf("%s:%d", conf.AppServer.Host, conf.AppServer.Port)

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      20 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	// Wait for signal or server error
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("received signal: %s, shutting down server...", sig)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v\n", err)
		} else {
			log.Println("server stopped gracefully")
		}
	}

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("graceful shutdown failed")
	} else {
		log.Println("server stopped gracefully")
	}
}
