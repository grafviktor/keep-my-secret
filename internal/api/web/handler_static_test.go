package web

/*
import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/grafviktor/keep-my-secret/internal/api"
)

func TestRegisterStaticHandler(t *testing.T) {
	// Create a new chi router
	router := chi.NewRouter()

	apiProvider := &apiRouteProvider{
		storage:  &MockStorage{},
		keyCache: &MockKeyCache{},
	}

	router.Get("/static/{id}", apiProvider.DownloadSecretFileHandler)

	// Create an HTTP request to test the router
	req := httptest.NewRequest("GET", "/static/valid_id", nil)
	req = req.WithContext(context.WithValue(req.Context(), api.ContextUserLogin, "validLogin"))
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Check the response status code (assuming the file exists)
	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.Code)
	}
}
*/
