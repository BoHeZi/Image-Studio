package main

import (
	"embed"
	"net/http"
	"runtime"

	"image-studio/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsmac "github.com/wailsapp/wails/v2/pkg/options/mac"
	wailswindows "github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	svc := backend.NewService()
	appOptions := &options.App{
		Frameless: true,
		Title:     "Image Studio",
		Width:     1440,
		Height:    980,
		MinWidth:  1100,
		MinHeight: 780,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Handler:    svc.MediaHandler(http.NotFoundHandler()),
			Middleware: svc.MediaHandler,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        svc.Startup,
		Bind: []interface{}{
			svc,
		},
	}

	if runtime.GOOS == "darwin" {
		appOptions.Mac = &wailsmac.Options{
			Appearance:           wailsmac.DefaultAppearance,
			TitleBar:             wailsmac.TitleBarHiddenInset(),
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		}
	}
	if runtime.GOOS == "windows" {
		webviewUserDataPath, err := backend.WindowsWebviewUserDataPath()
		if err != nil {
			println("Error:", err.Error())
			return
		}
		legacyWebviewUserDataPaths, err := backend.WindowsLegacyWebviewUserDataPaths()
		if err != nil {
			println("Error:", err.Error())
			return
		}
		if err := backend.MigrateWindowsWebviewDataDirs(webviewUserDataPath, legacyWebviewUserDataPaths); err != nil {
			println("Warning:", err.Error())
		}
		appOptions.Windows = &wailswindows.Options{
			Theme:                wailswindows.SystemDefault,
			BackdropType:         wailswindows.Mica,
			WebviewIsTransparent: false,
			WindowIsTranslucent:  true,
			WebviewUserDataPath:  webviewUserDataPath,
			CustomTheme: &wailswindows.ThemeSettings{
				DarkModeBorder:          wailswindows.RGB(54, 54, 54),
				DarkModeBorderInactive:  wailswindows.RGB(45, 45, 45),
				LightModeBorder:         wailswindows.RGB(219, 219, 219),
				LightModeBorderInactive: wailswindows.RGB(226, 226, 226),
			},
		}
	}

	err := wails.Run(appOptions)

	if err != nil {
		println("Error:", err.Error())
	}
}
