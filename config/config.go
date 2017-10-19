package config

import (
	"fmt"
	"os"
	"runtime"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var Conf Config

type Config struct {
	Core CoreConfig
	Jira JiraConfig
	Github GithubConfig
}

type CoreConfig struct {
	Editor string `toml:"editor"`
	SelectCmd string `toml:"selectcmd"`
}

type JiraConfig struct {
	UserName string `toml:"username"`
	Password string `toml:"password"`
	Jql string `toml:"jql"`
}

type GithubConfig struct {
	Template string `toml:"template"`
}

func (conf *Config) Load(file string) error {
	_, err := os.Stat(file)

	if err == nil {
		_, err := toml.DecodeFile(file, conf)
		if err != nil {
			return err
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	f, err := os.Create(file)

	if err != nil {
		return err
	}

	conf.Core.Editor = os.Getenv("EDITOR")
	if conf.Core.Editor == "" && runtime.GOOS != "windows" {
		conf.Core.Editor = "vim"
	}
	conf.Core.SelectCmd = "fzf"

	return toml.NewEncoder(f).Encode(conf)
}


func GetDefaultConfigDir() (dir string, err error) {
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "brant")
		}
		dir = filepath.Join(dir, "pet")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "brant")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("cannot create directory: %v", err)
	}
	return dir, nil
}
