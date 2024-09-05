package main

import (
	"net/http"

	"github.com/justinas/alice" // New import
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", app.health)
	mux.Handle("GET /metrics", app.promHandler())

	standard := alice.New(
		app.recoverPanic,
		app.logRequest,
		commonHeaders,
		app.getMetricsFromStatus,
		app.getMetricsFromNetInfo,
	)

	return standard.Then(mux)
}
