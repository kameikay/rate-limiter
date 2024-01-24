package webserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewWebServer(t *testing.T) {
	webserver := NewWebServer(":3333")
	assert.NotNil(t, webserver)
}

func TestMountMiddlewares(t *testing.T) {
	webserver := NewWebServer(":3333")
	webserver.MountMiddlewares()

	router := chi.NewRouter()
	router.Use(webserver.Router.Middlewares()...)

	// Define a test route to test the middlewares
	router.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a request to the test route
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a recorder to record the response
	rec := httptest.NewRecorder()

	// Call the test route with the request and recorder
	router.ServeHTTP(rec, req)

	// Assert that the response has the expected status code
	assert.Equal(t, http.StatusOK, rec.Code)
}
