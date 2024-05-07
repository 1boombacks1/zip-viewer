package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/1boombacks1/zipViewer/internal/model"
	"github.com/1boombacks1/zipViewer/internal/router"
	"github.com/1boombacks1/zipViewer/internal/zip"
)

type ZipReader interface {
	GetNamesByExt(string, string) ([]model.File, error)
}

type Application struct {
	port      int
	pathToZip string
	ext       string

	zipReader ZipReader
}

func New(path string) *Application {
	app := &Application{
		port:      8080,
		pathToZip: path,
		ext:       "*",

		zipReader: zip.New(),
	}

	return app
}

func (a *Application) SetPort(port int) {
	a.port = port
}

func (a *Application) SetExt(ext string) {
	ext = strings.TrimPrefix(ext, ".")
	a.ext = ext
}

func (app *Application) Start() {
	slog.Info("Getting names from zip with the specified extension", "path", app.pathToZip, "ext", app.ext)
	names, err := app.zipReader.GetNamesByExt(app.pathToZip, app.ext)
	if err != nil {
		slog.Error("Error getting names from zip", "path", app.pathToZip, "ext", app.ext, "error", err)
		return
	}
	slog.Info("A list of names has been received", "length", len(names))

	serverNotify := make(chan error)
	server := &http.Server{Addr: fmt.Sprintf(":%d", app.port), Handler: router.New(names, app.ext)}
	go func() {
		slog.Info("Listening on port", "port", app.port)
		serverNotify <- server.ListenAndServe()
		close(serverNotify)
	}()

	// Waiting signal
	slog.Info("Configuring graceful shutdown...")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-quit:
		slog.Info("app - Start - signal: " + s.String())
	case err = <-serverNotify:
		slog.Error("The server crashed with an error", "error", fmt.Errorf("app - Start - server notify: %w", err))
	}

	// Graceful shutdown
	slog.Info("Gracefully shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		slog.Error("Error at gracefully shutdown server", fmt.Errorf("app - Start - server.Shutdown: %w", err))
	}

}
