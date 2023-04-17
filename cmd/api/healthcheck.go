package main

import (
	"net/http"
)

type healthcheck struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

// Declare a handler which writes a json response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	h := healthcheck{
		Status:      "available",
		Environment: app.config.env,
		Version:     version,
	}

	err := app.writeJSON(w, http.StatusOK, h, nil)
	if err != nil {
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
