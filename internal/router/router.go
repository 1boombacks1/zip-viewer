package router

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/1boombacks1/zipViewer/internal/model"
	"github.com/go-chi/chi/v5"
)

func New(fileNames []model.File, ext string) http.Handler {
	r := chi.NewRouter()

	r.HandleFunc("/", showFiles(fileNames, ext))

	return r
}

func showFiles(files []model.File, ext string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("accept request", "addr", r.RemoteAddr, "path", r.URL.Path)

		templates := []string{
			"./ui/templates/main.layout.tmpl",
		}
		ts, err := template.ParseFiles(templates...)
		if err != nil {
			slog.Error("showFiles - template.ParseFiles: ", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		type templateData struct {
			Files     []model.File
			FileCount int
			Ext       string
		}
		data := &templateData{Files: files, FileCount: len(files), Ext: ext}

		if err := ts.Execute(w, data); err != nil {
			slog.Error("showFiles - template.Execute: ", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}

}
