package main

import (
	"context"
	"encoding/base64"
	"os"
	"strings"
	"path/filepath"
	"github.com/wailsapp/wails/v2/pkg/runtime"
    "fmt"
)

// App struct
type App struct {
	ctx context.Context
	lastDirectory string
	lastFilePath  string
  dirty         bool  
  allowQuitAnyway bool
}

func (a *App) SetDirty() {
	a.dirty = true
}

func (a *App) ClearDirty() {
	a.dirty = false
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// LoadFile opens a .ssg file and returns its contents
func (a *App) LoadFile() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Open a .ssg file",
		DefaultDirectory: a.lastDirectory,
		Filters: []runtime.FileFilter{
			{DisplayName: "Systemigram DSL files", Pattern: "*.ssg"},
		},
	})
	if err != nil || path == "" {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	a.lastDirectory = filepath.Dir(path)
	a.lastFilePath = path
	return string(data), nil
}


// SaveFile writes content to a .ssg file chosen by the user
func (a *App) SaveFile(content string) error {
	if a.lastFilePath != "" {
		return os.WriteFile(a.lastFilePath, []byte(content), 0644)
	}
	return a.SaveAsFile(content)
}

// Save As
func (a *App) SaveAsFile(content string) error {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "Save .ssg file as",
		DefaultFilename:  "diagram.ssg",
		DefaultDirectory: a.lastDirectory,
		Filters: []runtime.FileFilter{
			{DisplayName: "Systemigram DSL files", Pattern: "*.ssg"},
		},
	})
    if err != nil {
        return err
    }
    if path == "" {
        return fmt.Errorf("SaveAs canceled by user")
    }
	a.lastDirectory = filepath.Dir(path)
	a.lastFilePath = path
	return os.WriteFile(path, []byte(content), 0644)
}

func (a *App) SavePNGFile(base64data string) error {
	// Strip the prefix: "data:image/png;base64,"
	if idx := strings.Index(base64data, ","); idx != -1 {
		base64data = base64data[idx+1:]
	}

	// Decode base64
	imageData, err := base64.StdEncoding.DecodeString(base64data)
	if err != nil {
		return err
	}

	// Ask for save path
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save diagram as PNG",
	    DefaultDirectory: a.lastDirectory,
		DefaultFilename: "diagram.png",
		Filters: []runtime.FileFilter{
			{DisplayName: "PNG Image", Pattern: "*.png"},
		},
	})
	if err != nil || path == "" {
		return err
	}

	// Write file
    a.lastDirectory = filepath.Dir(path)
	return os.WriteFile(path, imageData, 0644)
}

// SaveDotFile writes the provided DOT content to a .dot file
func (a *App) SaveDotFile(content string) error {
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save Graphviz DOT file",
    	DefaultDirectory: a.lastDirectory,
		DefaultFilename: "systemigram.dot",
		Filters: []runtime.FileFilter{
			{DisplayName: "Graphviz DOT files", Pattern: "*.dot"},
		},
	})
	if err != nil || path == "" {
		return err
	}
    a.lastDirectory = filepath.Dir(path)
	return os.WriteFile(path, []byte(content), 0644)
}

func (a *App) OpenURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

func (a *App) QuitApp() {
	runtime.Quit(a.ctx)
}

func (a *App) AllowQuitAnyway() {
	a.allowQuitAnyway = true
}

func (a *App) OnBeforeClose(ctx context.Context) (prevent bool) {
	if a.allowQuitAnyway {
		return false // allow quit
	}
	if a.dirty {
		runtime.EventsEmit(ctx, "frontend-request-close")
		return true // prevent quit
	}
	return false
}

func (a *App) ForceQuit() {
	println("ForceQuit called")
	a.allowQuitAnyway = true
	runtime.Quit(a.ctx)
}

func (a *App) ClearLastFile() {
    a.lastFilePath = ""
}