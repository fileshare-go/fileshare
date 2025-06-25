package main

import (
	"context"
	"fmt"

	"github.com/chanmaoganda/fileshare/internal/fileshare"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var filters = []runtime.FileFilter{
	{
		DisplayName: "All Files (*.*)",
		Pattern:     "*.*",
	},
	{
		DisplayName: "Text Files (*.txt)",
		Pattern:     "*.txt",
	},
	{
		DisplayName: "Images (*.png;*.jpg)",
		Pattern:     "*.png;*.jpg",
	},
}

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) OpenFileSelectionDialog() (string, error) {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Select a File",
		Filters:              filters,
		CanCreateDirectories: true,
		ShowHiddenFiles:      false,
	})

	if err != nil {
		fmt.Printf("File dialog error: %v\n", err)
		return "", err
	}

	return selection, nil
}

func (a *App) UploadFile(args []string) {
	fileshare.Upload(context.TODO(), args)
}
