package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/grafviktor/keep-my-secret/internal/config"
)

func TestRegisterStaticHandler(t *testing.T) {
	// Create a new chi router for testing.
	serveDir = "./ui"
	r := chi.NewRouter()

	webAppUrls := []string{"/valid/", "/missing_ending_slash", "" /*empty*/}

	for _, url := range webAppUrls {
		// Define a test AppConfig.
		testConfig := config.AppConfig{
			ClientAppURL: url,
		}

		// Register the static handler.
		registerStaticHandler(testConfig, r)

		// Create a test request for the root URL.
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		// Execute the request.
		r.ServeHTTP(w, req)

		if !strings.HasSuffix(url, "/") {
			url += "/"
		}

		// I know that this file exists in project_root/ui/.keep
		filePath := url + ".keep"
		req = httptest.NewRequest("GET", filePath, nil)
		w = httptest.NewRecorder()

		// Execute the request.
		r.ServeHTTP(w, req)

		// Check the response status code. It should be StatusOK if the file exists.
		assert.Equal(t, http.StatusOK, w.Code)

		// Create a test request for a non-existent file.
		nonExistentFilePath := "/static/non_existent_file.txt"
		req = httptest.NewRequest("GET", nonExistentFilePath, nil)
		w = httptest.NewRecorder()

		// Execute the request.
		r.ServeHTTP(w, req)

		// Check the response status code. It should be StatusNotFound.
		assert.Equal(t, http.StatusNotFound, w.Code)
	}
}

func init() {
	projectRoot := "../../../"
	err := os.Chdir(projectRoot)
	if err != nil {
		panic(err)
	}
}
