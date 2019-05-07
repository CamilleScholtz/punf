package main

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
)

// config is a stuct with all config values. See `runtime/config/config.toml`
// for more information about these values.
var config struct {
	User string
	Pass string

	Scrot    string
	SelScrot string

	Clipboard bool
	Print     bool
}

// parseConfig parses a toml config.
func parseConfig() error {
	hd, err := homedir.Dir()
	if err != nil {
		return err
	}

	if _, err = toml.DecodeFile(filepath.Join(hd, ".punf", "config.toml"),
		&config); err != nil {
		return fmt.Errorf("config %s: %s", filepath.Join(hd, ".punf",
			"config.toml"), err)
	}

	return nil
}
