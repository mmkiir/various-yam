package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

// App struct
type App struct {
	ctx context.Context
	fs  *FileStorage
}

// NewApp creates a new App application struct
func NewApp() *App {
	if err := os.MkdirAll(filepath.Join(xdg.DataHome, "various-yam"), 0755); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	file, err := os.OpenFile(
		filepath.Join(xdg.DataHome, "various-yam", "data.json"),
		os.O_CREATE,
		0644,
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	file.Close()

	return &App{
		fs: NewFileStorage(filepath.Join(xdg.DataHome, "various-yam", "data.json")),
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}
