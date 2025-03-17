package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/api/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/movies/:id", app.showMovieHandler)

	return router
}
