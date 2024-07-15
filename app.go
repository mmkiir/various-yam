package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/gen2brain/malgo"
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

// MediaDeviceInfo struct
type MediaDeviceInfo struct {
	DeviceID string `json:"deviceId"`
	GroupID  string `json:"groupId"`
	Kind     string `json:"kind"`
	Label    string `json:"label"`
}

// ListCaptureDevices lists all available capture devices
func (a *App) ListCaptureDevices() ([]MediaDeviceInfo, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		slog.Info(message)
	})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer func() {
		if err := ctx.Uninit(); err != nil {
			slog.Error(err.Error())
		}
		ctx.Free()
	}()

	deviceInfos, err := ctx.Devices(malgo.Capture)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	var devices []MediaDeviceInfo
	for _, deviceInfo := range deviceInfos {
		devices = append(devices, MediaDeviceInfo{
			DeviceID: deviceInfo.ID.String(),
			Kind:     "audioinput",
			Label:    deviceInfo.Name(),
		})
	}

	return devices, nil
}

// GetCaptureDevice gets the capture device by ID
func (a *App) GetCaptureDeviceID() (string, error) {
	serializedCaptureDeviceID, _ := a.fs.GetItem("captureDeviceID")
	if serializedCaptureDeviceID == "" {
		return "", nil
	}

	captureDeviceID := ""
	if err := json.Unmarshal([]byte(serializedCaptureDeviceID), &captureDeviceID); err != nil {
		return "", err
	}

	return captureDeviceID, nil
}

// SetCaptureDeviceID sets the capture device by ID
func (a *App) SetCaptureDeviceID(captureDeviceID string) error {
	serializedCaptureDeviceID, err := json.Marshal(captureDeviceID)
	if err != nil {
		return err
	}

	if err := a.fs.SetItem("captureDeviceID", string(serializedCaptureDeviceID)); err != nil {
		return err
	}

	return nil
}

// ListPlaybackDevices lists all available playback devices
func (a *App) ListPlaybackDevices() ([]MediaDeviceInfo, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(message string) {
		slog.Info(message)
	})
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer func() {
		if err := ctx.Uninit(); err != nil {
			slog.Error(err.Error())
		}
		ctx.Free()
	}()

	deviceInfos, err := ctx.Devices(malgo.Playback)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	var devices []MediaDeviceInfo
	for _, deviceInfo := range deviceInfos {
		devices = append(devices, MediaDeviceInfo{
			DeviceID: deviceInfo.ID.String(),
			Kind:     "audiooutput",
			Label:    deviceInfo.Name(),
		})
	}

	return devices, nil
}

// GetPlaybackDeviceID gets the playback device by ID
func (a *App) GetPlaybackDeviceID() (string, error) {
	serializedPlaybackDeviceID, _ := a.fs.GetItem("playbackDeviceID")
	if serializedPlaybackDeviceID == "" {
		return "", nil
	}

	playbackDeviceID := ""
	if err := json.Unmarshal([]byte(serializedPlaybackDeviceID), &playbackDeviceID); err != nil {
		return "", err
	}

	return playbackDeviceID, nil
}

// SetPlaybackDeviceID sets the playback device by ID
func (a *App) SetPlaybackDeviceID(playbackDeviceID string) error {
	serializedPlaybackDeviceID, err := json.Marshal(playbackDeviceID)
	if err != nil {
		return err
	}

	if err := a.fs.SetItem("playbackDeviceID", string(serializedPlaybackDeviceID)); err != nil {
		return err
	}

	return nil
}
