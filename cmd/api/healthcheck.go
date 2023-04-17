package main

import (
	"encoding/json"
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

	// Pass the struct to the json.Marshal() function. This returns a []byte slice
	// containing the encoded JSON. If there was an error, we log it and send the client
	// a generic error message.
	js, err := json.Marshal(h)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)

		return
	}

	// Set the "Content-Type: application/json" header on the response. If you forget to
	// this, Go will default to sending a "Content-Type: text/plain; charset=utf-8"
	// header instead.
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON as the HTTP response body.
	// What do we do if there is an error? log it.
	// See https://stackoverflow.com/a/43976633/2180697 for more info.
	if _, err := w.Write(js); err != nil {
		app.logger.Println("an error occurred writing response: %w", err)
	}
}
