package main

import (
	"net/http"
)

func (app *Application) health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
