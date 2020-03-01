package main

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

// config is a stuct with all config values. See `runtime/config/config.toml`
// for more information about these values.
var config struct {
	URL string

	User string
	Pass string

	Scrot    string
	SelScrot string

	Clipboard bool
	Print     bool
}

// parseConfig parses a toml config.
func parseConfig() error {
	if _, err := toml.DecodeFile(filepath.Join(xdg.ConfigHome, "punf",
		"config.toml"), &config); err != nil {
		return fmt.Errorf("config %s: %s", filepath.Join(xdg.ConfigHome, "punf",
			"config.toml"), err)
	}

	return nil
}
