package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"goservice/configs"
	"goservice/internal/client"
	"goservice/internal/student"
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

	r.Mount("/students", handler.Routes())

	addr := fmt.Sprintf(":%d", conf.AppServer.Port)

	fmt.Println("🚀 Server running on :", addr)
	http.ListenAndServe(addr, r)
}
