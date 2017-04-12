package main

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	homedir "github.com/mitchellh/go-homedir"
)

// config is a stuct with all config values. See `runtime/config/config.toml` for
// more information about these values.
var config struct {
	Host string
	Key  string

	Scrot    string
	SelScrot string

	Clipboard    bool
	Log          bool
	Notification bool
	Print        bool
}

// initConfig initializes the config struct.
func initConfig() error {
	hd, err := homedir.Dir()
	if err != nil {
		return err
	}

	if _, err = toml.DecodeFile(filepath.Join(hd, "punf/config.toml"), &config); err != nil {
		return fmt.Errorf("config ./config.toml: " + err.Error())
	}

	return nil
}
