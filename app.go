package main

import (
	"controlplane/handler"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
}

func (app *App) Run(port string) {
	app.Build()
	log.Fatal(http.ListenAndServe(port, app.Router))
}

func (app *App) Build() {
	app.Router = mux.NewRouter()
	app.Router.HandleFunc("/api", handler.GetAllApi).Methods("GET")
	app.Router.HandleFunc("/api/{id}", handler.GetApi).Methods("GET")
	app.Router.HandleFunc("/api", handler.CreateApi).Methods("POST")
	app.Router.HandleFunc("/api/{id}/spec/{version}", handler.UploadSpec).Methods("POST")
	app.Router.HandleFunc("/api/{id}", handler.UpdateApi).Methods("PUT")
	app.Router.HandleFunc("/api/{id}", handler.DeleteApi).Methods("DELETE")
	app.Router.HandleFunc("/api/{api_id}/route", handler.CreateRoute).Methods("POST")
	app.Router.HandleFunc("/api/{api_id}/route/{route_id}", handler.UpdateRoute).Methods("PUT")
	app.Router.HandleFunc("/api/{api_id}/route/{route_id}", handler.GetRoute).Methods("GET")
}
