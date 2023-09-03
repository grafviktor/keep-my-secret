package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafviktor/keep-my-secret/internal/api"
	"github.com/grafviktor/keep-my-secret/internal/constant"
	"github.com/grafviktor/keep-my-secret/internal/version"
)

func TestVersionHandler(t *testing.T) {
	// Create a mock HTTP request
	req := httptest.NewRequest("GET", "/version", nil)

	// Create a mock HTTP response recorder
	res := httptest.NewRecorder()

	// Call the VersionHandler function to handle the request
	VersionHandler(res, req)

	// Check the response status code
	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, res.Code)
	}

	// Parse the response JSON
	var response api.Response
	err := json.Unmarshal(res.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response JSON: %v", err)
	}

	// Check the status field in the response
	if response.Status != constant.APIStatusSuccess {
		t.Errorf("Expected response status '%s', but got '%s'", constant.APIStatusSuccess, response.Status)
	}

	// Check the data field in the response
	var versionResp versionResponse
	jsonResp, err := json.Marshal(response.Data)
	if err != nil {
		t.Fatalf("Failed to marshal response.data into json: %v", err)
	}

	err = json.Unmarshal(jsonResp, &versionResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal version response: %v", err)
	}

	// Validate the fields in the version response
	expectedVersion := version.BuildVersion()
	expectedDate := version.BuildDate()
	expectedCommit := version.BuildCommit()
	expectedAPIVersion := "1.0.0"

	if versionResp.BuildVersion != expectedVersion {
		t.Errorf("Expected build version '%s', but got '%s'", expectedVersion, versionResp.BuildVersion)
	}

	if versionResp.BuildDate != expectedDate {
		t.Errorf("Expected build date '%s', but got '%s'", expectedDate, versionResp.BuildDate)
	}

	if versionResp.BuildCommit != expectedCommit {
		t.Errorf("Expected build commit '%s', but got '%s'", expectedCommit, versionResp.BuildCommit)
	}

	if versionResp.APIVersion != expectedAPIVersion {
		t.Errorf("Expected API version '%s', but got '%s'", expectedAPIVersion, versionResp.APIVersion)
	}
}
