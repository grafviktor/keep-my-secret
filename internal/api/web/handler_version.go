package web

import (
	"github.com/grafviktor/keep-my-secret/internal/api"
	"github.com/grafviktor/keep-my-secret/internal/api/utils"
	"github.com/grafviktor/keep-my-secret/internal/constant"
	"github.com/grafviktor/keep-my-secret/internal/version"
	"net/http"
)

type versionResponse struct {
	BuildVersion string `json:"build_version"`
	BuildDate    string `json:"build_date"`
	BuildCommit  string `json:"build_commit"`
	ApiVersion   string `json:"api_version"`
}

func VersionHandler(w http.ResponseWriter, _ *http.Request) {
	_ = utils.WriteJSON(w, http.StatusOK, api.Response{
		Status: constant.APIStatusSuccess,
		Data: versionResponse{
			BuildVersion: version.BuildVersion(),
			BuildDate:    version.BuildDate(),
			BuildCommit:  version.BuildCommit(),
			ApiVersion:   "1.0.0",
		},
	})
}
