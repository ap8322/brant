package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

var Conf Config

type Config struct {
	Core CoreConfig
	Jira JiraConfig
}

type CoreConfig struct {
	Editor      string `toml:"editor"`
	TicketCache string `toml:"ticketcache"`
	SelectCmd   string `toml:"selectcmd"`
	baseBranch  string `toml:"basebranch"`
}

type JiraConfig struct {
	Host     string `toml:"host"`
	UserName string `toml:"username"`
	Password string `toml:"password"`
	Jql      string `toml:"jql"`
}

func (conf *Config) Load(file string) error {
	_, err := os.Stat(file)

	if err == nil {
		_, err := toml.DecodeFile(file, conf)
		if err != nil {
			return err
		}
		conf.Core.TicketCache = expandPath(conf.Core.TicketCache)
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	dir, _ := GetDefaultConfigDir()
	conf.Core.TicketCache = filepath.Join(dir, "tickets.toml")
	_, err = os.Create(conf.Core.TicketCache)
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
		dir = filepath.Join(dir, "brant")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "brant")
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("cannot create directory: %v", err)
	}
	return dir, nil
}

func expandPath(s string) string {
	if len(s) >= 2 && s[0] == '~' && os.IsPathSeparator(s[1]) {
		if runtime.GOOS == "windows" {
			s = filepath.Join(os.Getenv("USERPROFILE"), s[2:])
		} else {
			s = filepath.Join(os.Getenv("HOME"), s[2:])
		}
	}
	return os.Expand(s, os.Getenv)
}
