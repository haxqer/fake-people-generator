package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/index.html static/styles.css static/app.js
var embeddedFiles embed.FS

var (
	staticFiles fs.FS
	indexHTML   []byte
)

func init() {
	var err error

	staticFiles, err = fs.Sub(embeddedFiles, "static")
	if err != nil {
		panic(err)
	}

	indexHTML, err = embeddedFiles.ReadFile("static/index.html")
	if err != nil {
		panic(err)
	}
}

func IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(indexHTML)
	}
}

func AssetsHandler() http.Handler {
	return http.FileServerFS(staticFiles)
}
