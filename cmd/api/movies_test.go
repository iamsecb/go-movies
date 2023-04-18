package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	// Create an instance of the application. At this time it can be reused across our tests.
	app application
	// Generate the http routes and save the handler for reuse across all our tests.
	handler = http.HandlerFunc(app.routes().ServeHTTP)
)

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler.
	req, err := http.NewRequest(http.MethodGet, "/api/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call our ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	want := `{"status":"available","environment":%q,"version":%q}`
	want = fmt.Sprintf(want, app.config.env, version)

	if got := strings.TrimSpace(rr.Body.String()); got != want {
		t.Errorf("handler returned wrong response: got %v want %v", got, want)
	}
}

func TestCreateMovieHandler(t *testing.T) {
	testData := []struct {
		name               string
		expectedStatusCode int
		body               string
	}{
		{name: "success",
			body: `{
			"title":"Moana",
			"year":2016,
			"runtime":"107 mins",
			"genres":["animation","adventure"]
		   }
		   `,
			expectedStatusCode: 200,
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest(http.MethodPost, "/api/v1/movies", strings.NewReader(td.body))
			if err != nil {
				t.Fatal(err)
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != td.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, td.expectedStatusCode)
			}
		})
	}
}

func TestShowMovieHandler(t *testing.T) {
	testData := []struct {
		name               string
		id                 string
		expectedStatusCode int
	}{
		{name: "valid", id: "1", expectedStatusCode: 200},
		{name: "invalid", id: "blah", expectedStatusCode: 404},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
			// pass 'nil' as the third parameter.
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/movies/%s", td.id), nil)
			if err != nil {
				t.Fatal(err)
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != td.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, td.expectedStatusCode)
			}
		})
	}
}
