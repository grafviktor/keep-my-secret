package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/grafviktor/keep-my-secret/internal/config"
)

var serveDir = "./ui/static"

// registerStaticHandler - registers a handler for serving client application.
func registerStaticHandler(config config.AppConfig, router *chi.Mux) {
	workDir, _ := os.Getwd()
	fsPath := http.Dir(filepath.Join(workDir, serveDir))
	url := config.ClientAppURL

	if len(url) == 0 {
		url = "/"
	}

	if strings.ContainsAny(url, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if url != "/" && url[len(url)-1] != '/' {
		url += "/"
	}
	url += "*"

	router.Get(url, func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(fsPath))
		fs.ServeHTTP(w, r)
	})
}
