package api

import (
	"embed"
	"io"
	"net/http"
)

//go:embed web/index.html
var webFiles embed.FS

func serveHome(w http.ResponseWriter, _ *http.Request) {
	page, err := webFiles.Open("web/index.html")
	if err != nil {
		http.Error(w, "page unavailable", http.StatusInternalServerError)
		return
	}
	defer page.Close()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, page)
}
