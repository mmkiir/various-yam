package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/gen2brain/malgo"
	"github.com/hajimehoshi/go-mp3"
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/youpy/go-wav"
)

// App struct
type App struct {
	ctx                          context.Context
	fs                           *FileStorage
	cancelFunctionsForAudioFiles map[string]context.CancelFunc
	cancelLoopbackAudio          context.CancelFunc
}

// NewApp creates a new App application struct
func NewApp() *App {
	dataPath := filepath.Join(xdg.DataHome, "various-yam")
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		log.Fatal("Failed to create data directory: ", err)
	}

	filePath := filepath.Join(dataPath, "data.json")
	file, err := os.OpenFile(filePath, os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("Failed to open data file: ", err)
	}
	file.Close()

	return &App{
		fs:                           NewFileStorage(filePath),
		cancelFunctionsForAudioFiles: make(map[string]context.CancelFunc),
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	ctx, cancel := context.WithCancel(ctx)
	a.cancelLoopbackAudio = cancel

	go func() {
		if err := a.LoopbackAudio(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	go a.registerAudioFileKeybindings()
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit
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
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := ctx.Uninit(); err != nil {
			log.Fatal(err)
		}
		ctx.Free()
	}()

	deviceInfos, err := ctx.Devices(malgo.Capture)
	if err != nil {
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

// GetCaptureDeviceID gets the capture device by ID
func (a *App) GetCaptureDeviceID() (string, error) {
	serializedCaptureDeviceID, _ := a.fs.GetItem("captureDeviceID")
	if serializedCaptureDeviceID == "" {
		return "", nil
	}

	var captureDeviceID string
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

	a.cancelLoopbackAudio()

	ctx, cancel := context.WithCancel(a.ctx)
	a.cancelLoopbackAudio = cancel

	go func() {
		if err := a.LoopbackAudio(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

// ListPlaybackDevices lists all available playback devices
func (a *App) ListPlaybackDevices() ([]MediaDeviceInfo, error) {
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := ctx.Uninit(); err != nil {
			log.Fatal(err)
		}
		ctx.Free()
	}()

	deviceInfos, err := ctx.Devices(malgo.Playback)
	if err != nil {
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

	var playbackDeviceID string
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

	a.cancelLoopbackAudio()

	ctx, cancel := context.WithCancel(a.ctx)
	a.cancelLoopbackAudio = cancel

	go func() {
		if err := a.LoopbackAudio(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

// ListAudioFiles lists all available audio files
func (a *App) ListAudioFiles() ([]string, error) {
	serializedAudioFiles, _ := a.fs.GetItem("audioFiles")
	if serializedAudioFiles == "" {
		return []string{}, nil
	}

	var audioFiles []string
	if err := json.Unmarshal([]byte(serializedAudioFiles), &audioFiles); err != nil {
		return nil, err
	}

	return audioFiles, nil
}

// AddAudioFile adds an audio file
func (a *App) AddAudioFile(audioFile string) error {
	return a.fs.UpdateItem("audioFiles", func(value string) (string, error) {
		var audioFiles []string
		if value != "" {
			if err := json.Unmarshal([]byte(value), &audioFiles); err != nil {
				return "", err
			}
		}

		audioFiles = append(audioFiles, audioFile)

		serializedAudioFiles, err := json.Marshal(audioFiles)
		if err != nil {
			return "", err
		}

		return string(serializedAudioFiles), nil
	})
}

// RemoveAudioFile removes an audio file
func (a *App) RemoveAudioFile(audioFile string) error {
	return a.fs.UpdateItem("audioFiles", func(value string) (string, error) {
		var audioFiles []string
		if value != "" {
			if err := json.Unmarshal([]byte(value), &audioFiles); err != nil {
				return "", err
			}
		}

		for i, file := range audioFiles {
			if file == audioFile {
				audioFiles = append(audioFiles[:i], audioFiles[i+1:]...)
				break
			}
		}

		serializedAudioFiles, err := json.Marshal(audioFiles)
		if err != nil {
			return "", err
		}

		return string(serializedAudioFiles), nil
	})
}

// PlayAudioFile plays an audio file
func (a *App) PlayAudioFile(audioFile string) error {
	a.StopAudioFile(audioFile)

	ctx, cancel := context.WithCancel(a.ctx)
	a.cancelFunctionsForAudioFiles[audioFile] = cancel

	go func() {
		audioContext, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := audioContext.Uninit(); err != nil {
				log.Fatal(err)
			}
			audioContext.Free()
		}()

		file, err := os.Open(audioFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		var audioFormat malgo.FormatType
		var channels uint32
		var reader io.Reader
		var sampleRate uint32

		switch filepath.Ext(audioFile) {
		case ".wav":
			w := wav.NewReader(file)
			f, err := w.Format()
			if err != nil {
				log.Fatal(err)
			}

			switch f.AudioFormat {
			case 1:
				switch f.BitsPerSample {
				case 8:
					audioFormat = malgo.FormatU8
				case 16:
					audioFormat = malgo.FormatS16
				case 24:
					audioFormat = malgo.FormatS24
				case 32:
					audioFormat = malgo.FormatS32
				default:
					log.Fatal("Unsupported bits per sample: ", f.BitsPerSample)
				}
			case 3:
				switch f.BitsPerSample {
				case 32:
					audioFormat = malgo.FormatF32
				default:
					log.Fatal("Unsupported bits per sample: ", f.BitsPerSample)
				}
			default:
				log.Fatal("Unsupported audio format: ", f.AudioFormat)
			}

			channels = uint32(f.NumChannels)
			reader = w
			sampleRate = f.SampleRate
		case ".mp3":
			m, err := mp3.NewDecoder(file)
			if err != nil {
				log.Fatal(err)
			}

			audioFormat = malgo.FormatS16
			channels = 2
			reader = m
			sampleRate = uint32(m.SampleRate())
		default:
			log.Fatal("Unsupported audio file format: ", filepath.Ext(audioFile))
		}

		deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
		deviceConfig.Alsa.NoMMap = 1
		deviceConfig.Playback.Channels = channels
		deviceConfig.Playback.Format = audioFormat
		deviceConfig.SampleRate = sampleRate

		deviceID, _ := a.GetPlaybackDeviceID()
		if deviceID != "" {
			deviceID, err := ParseHexStringToDeviceID(deviceID)
			if err != nil {
				log.Fatal(err)
			}
			deviceConfig.Playback.DeviceID = deviceID.Pointer()
		}

		deviceCallbacks := malgo.DeviceCallbacks{
			Data: func(pOutputSample, _ []byte, _ uint32) {
				io.ReadFull(reader, pOutputSample)
			},
		}

		device, err := malgo.InitDevice(audioContext.Context, deviceConfig, deviceCallbacks)
		if err != nil {
			log.Fatal(err)
		}

		if err := device.Start(); err != nil {
			log.Fatal(err)
		}

		<-ctx.Done()

		if err := device.Stop(); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

// StopAudioFile stops an audio file
func (a *App) StopAudioFile(audioFile string) error {
	cancel, ok := a.cancelFunctionsForAudioFiles[audioFile]
	if !ok {
		return nil
	}

	cancel()

	delete(a.cancelFunctionsForAudioFiles, audioFile)

	return nil
}

// FileFilter defines a filter for dialog boxes
type FileFilter struct {
	DisplayName string `json:"displayName"` // Filter information EG: "Image Files (*.jpg, *.png)"
	Pattern     string `json:"pattern"`     // semicolon separated list of extensions, EG: "*.jpg;*.png"
}

// OpenDialogOptions contains the options for the OpenDialogOptions runtime method
type OpenDialogOptions struct {
	DefaultDirectory           string       `json:"defaultDirectory"`
	DefaultFilename            string       `json:"defaultFilename"`
	Title                      string       `json:"title"`
	Filters                    []FileFilter `json:"filters"`
	ShowHiddenFiles            bool         `json:"showHiddenFiles"`
	CanCreateDirectories       bool         `json:"canCreateDirectories"`
	ResolvesAliases            bool         `json:"resolvesAliases"`
	TreatPackagesAsDirectories bool         `json:"treatPackagesAsDirectories"`
}

// OpenMultipleFilesDialog prompts the user to select a file
func (a *App) OpenMultipleFilesDialog(dialogOptions OpenDialogOptions) ([]string, error) {
	filters := make([]runtime.FileFilter, len(dialogOptions.Filters))
	for i, filter := range dialogOptions.Filters {
		filters[i] = runtime.FileFilter{
			DisplayName: filter.DisplayName,
			Pattern:     filter.Pattern,
		}
	}

	return runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		DefaultDirectory:           dialogOptions.DefaultDirectory,
		DefaultFilename:            dialogOptions.DefaultFilename,
		Title:                      dialogOptions.Title,
		Filters:                    filters,
		ShowHiddenFiles:            dialogOptions.ShowHiddenFiles,
		CanCreateDirectories:       dialogOptions.CanCreateDirectories,
		ResolvesAliases:            dialogOptions.ResolvesAliases,
		TreatPackagesAsDirectories: dialogOptions.TreatPackagesAsDirectories,
	})
}

// LoopbackAudio loops back audio from the capture device to the playback device
func (a *App) LoopbackAudio(ctx context.Context) error {
	audioContext, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return err
	}

	captureDeviceConfig := malgo.DefaultDeviceConfig(malgo.Capture)
	captureDeviceConfig.Alsa.NoMMap = 1
	captureDeviceConfig.Capture.Channels = 1
	captureDeviceConfig.Capture.Format = malgo.FormatS16
	captureDeviceConfig.SampleRate = 44100

	captureDeviceID, _ := a.GetCaptureDeviceID()
	if captureDeviceID != "" {
		deviceID, err := ParseHexStringToDeviceID(captureDeviceID)
		if err != nil {
			return err
		}
		captureDeviceConfig.Capture.DeviceID = deviceID.Pointer()
	}

	pInputSamplesChannel := make(chan []byte)

	captureDeviceCallbacks := malgo.DeviceCallbacks{
		Data: func(_, pInputSamples []byte, _ uint32) {
			pInputSamplesChannel <- pInputSamples
		},
	}

	captureDevice, err := malgo.InitDevice(audioContext.Context, captureDeviceConfig, captureDeviceCallbacks)
	if err != nil {
		return err
	}

	if err := captureDevice.Start(); err != nil {
		return err
	}

	playbackDeviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
	playbackDeviceConfig.Alsa.NoMMap = 1
	playbackDeviceConfig.Playback.Channels = 1
	playbackDeviceConfig.Playback.Format = malgo.FormatS16
	playbackDeviceConfig.SampleRate = 44100

	playbackDeviceID, _ := a.GetPlaybackDeviceID()
	if playbackDeviceID != "" {
		deviceID, err := ParseHexStringToDeviceID(playbackDeviceID)
		if err != nil {
			return err
		}
		playbackDeviceConfig.Playback.DeviceID = deviceID.Pointer()
	}

	playbackDeviceCallbacks := malgo.DeviceCallbacks{
		Data: func(pOutputSample, _ []byte, _ uint32) {
			pInputSamples, ok := <-pInputSamplesChannel
			if !ok {
				return
			}
			copy(pOutputSample, pInputSamples)
		},
	}

	playbackDevice, err := malgo.InitDevice(audioContext.Context, playbackDeviceConfig, playbackDeviceCallbacks)
	if err != nil {
		return err
	}

	if err := playbackDevice.Start(); err != nil {
		return err
	}

	<-ctx.Done()

	if err := captureDevice.Stop(); err != nil {
		return err
	}

	if err := playbackDevice.Stop(); err != nil {
		return err
	}

	return nil
}

// ListAudioFileKeybindings lists all available audio file keybindings
func (a *App) ListAudioFileKeybindings() (map[string]string, error) {
	serializedAudioFileKeybindings, _ := a.fs.GetItem("audioFileKeybindings")
	if serializedAudioFileKeybindings == "" {
		return map[string]string{}, nil
	}

	var audioFileKeybindings map[string]string
	if err := json.Unmarshal([]byte(serializedAudioFileKeybindings), &audioFileKeybindings); err != nil {
		return nil, err
	}

	return audioFileKeybindings, nil
}

// SetAudioFileKeybinding sets the keybinding for an audio file
func (a *App) SetAudioFileKeybinding(audioFile string, keybinding string) error {
	if err := a.fs.UpdateItem("audioFileKeybindings", func(value string) (string, error) {
		audioFileKeybindings := make(map[string]string)
		if value != "" {
			if err := json.Unmarshal([]byte(value), &audioFileKeybindings); err != nil {
				return "", err
			}
		}

		audioFileKeybindings[audioFile] = keybinding

		serializedAudioFileKeybindings, err := json.Marshal(audioFileKeybindings)
		if err != nil {
			return "", err
		}

		return string(serializedAudioFileKeybindings), nil
	}); err != nil {
		return err
	}

	go a.reloadAudioFileKeybindings()

	return nil
}

// RemoveAudioFileKeybinding removes the keybinding for an audio file
func (a *App) RemoveAudioFileKeybinding(audioFile string) error {
	if err := a.fs.UpdateItem("audioFileKeybindings", func(value string) (string, error) {
		audioFileKeybindings := make(map[string]string)
		if value != "" {
			if err := json.Unmarshal([]byte(value), &audioFileKeybindings); err != nil {
				return "", err
			}
		}

		delete(audioFileKeybindings, audioFile)

		serializedAudioFileKeybindings, err := json.Marshal(audioFileKeybindings)
		if err != nil {
			return "", err
		}

		return string(serializedAudioFileKeybindings), nil
	}); err != nil {
		return err
	}

	go a.reloadAudioFileKeybindings()

	return nil
}

func (a *App) registerAudioFileKeybindings() error {
	audioFileKeybindings, err := a.ListAudioFileKeybindings()
	if err != nil {
		return err
	}

	for audioFile, keybinding := range audioFileKeybindings {
		if keybinding == "" {
			continue
		}

		audioFile := audioFile

		hook.Register(hook.KeyDown, strings.Split(strings.ToLower(keybinding), " + "), func(e hook.Event) {
			if err := a.PlayAudioFile(audioFile); err != nil {
				log.Fatal(err)
			}
		})
	}

	s := hook.Start()
	<-hook.Process(s)

	return nil
}

func (a *App) reloadAudioFileKeybindings() error {
	hook.End()

	return a.registerAudioFileKeybindings()
}

// ParseHexStringToDeviceID parses a hex string to a malgo.DeviceID
func ParseHexStringToDeviceID(s string) (malgo.DeviceID, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return malgo.DeviceID{}, err
	}

	if len(bytes) > len(malgo.DeviceID{}) {
		return malgo.DeviceID{}, fmt.Errorf("malgo.DeviceID is too short for the given string")
	}

	var deviceID malgo.DeviceID
	copy(deviceID[:], bytes)

	return deviceID, nil
}
