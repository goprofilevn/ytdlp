package main

import (
	"embed"
	"fmt"
	"ytdlp/helpers/logrus"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// init logger
	logrus.NewLogrusLogger()
	logrus.InitLogrusLogger()
	appVersion := "1.0.1"
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:         fmt.Sprintf("Yt-DLP %s", appVersion),
		Width:         700,
		Height:        500,
		DisableResize: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 247, G: 249, B: 252, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
