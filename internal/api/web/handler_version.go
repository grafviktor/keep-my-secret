package web

import (
	"net/http"

	"github.com/grafviktor/keep-my-secret/internal/api"
	"github.com/grafviktor/keep-my-secret/internal/api/utils"
	"github.com/grafviktor/keep-my-secret/internal/constant"
	"github.com/grafviktor/keep-my-secret/internal/version"
)

type versionResponse struct {
	BuildVersion string `json:"build_version"`
	BuildDate    string `json:"build_date"`
	BuildCommit  string `json:"build_commit"`
	APIVersion   string `json:"api_version"`
}

// VersionHandler returns the version of the application, including API version
func VersionHandler(w http.ResponseWriter, _ *http.Request) {
	_ = utils.WriteJSON(w, http.StatusOK, api.Response{
		Status: constant.APIStatusSuccess,
		Data: versionResponse{
			BuildVersion: version.BuildVersion(),
			BuildDate:    version.BuildDate(),
			BuildCommit:  version.BuildCommit(),
			APIVersion:   "1.0.0",
		},
	})
}
