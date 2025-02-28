//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"os"
	"path"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
)

type config struct {
	Cleanup bool `yaml:"cleanup" mapstructure:"cleanup"`
}

var (
	app *portapps.App
	cfg *config
)

func init() {
	var err error

	// Default config
	cfg = &config{
		Cleanup: false,
	}

	// Init app
	if app, err = portapps.NewWithCfg("AnotherRedisDesktopManager-portable", "AnotherRedisDesktopManager", cfg); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "Another Redis Desktop Manager.exe")
	app.Args = []string{
		"--user-data-dir=" + app.DataPath,
	}

	// Cleanup on exit
	if cfg.Cleanup {
		defer func() {
			utl.Cleanup([]string{
				path.Join(os.Getenv("APPDATA"), "AnotherRedisDesktopManager"),
			})
		}()
	}

	configFile := utl.PathJoin(app.DataPath, "config.yaml")
	if !utl.Exists(configFile) {
		log.Info().Msg("Creating default config...")
		if err := utl.WriteToFile(configFile, `enableAutomaticUpdates: false`); err != nil {
			log.Error().Err(err).Msg("Cannot write default config")
		}
	}
	if err := utl.ReplaceByPrefix(configFile, "enableAutomaticUpdates:", "enableAutomaticUpdates: false"); err != nil {
		log.Fatal().Err(err).Msg("Cannot set enableAutomaticUpdates property")
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
