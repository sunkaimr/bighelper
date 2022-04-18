package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

const (
	DefaultCfgFileName = "bighelper.ini"
)

func LoadConfig(cfgPath string) (*ini.File, error) {
	configPath, err := findConfig(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("Fail to find config file: %v", err)
	}

	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("Fail to load config file: %v", err)
	}
	return cfg, nil
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func findConfig(cfgPath string) (string, error) {
	if cfgPath == "" {
		cfgPath = filepath.Join("./", DefaultCfgFileName)
	}

	if !fileExist(cfgPath) {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		index := strings.LastIndex(path, string(os.PathSeparator))
		cfgPathAbs := filepath.Join(path[:index], DefaultCfgFileName)

		if !fileExist(cfgPathAbs) {
			return "", fmt.Errorf("config file not found in %v or %v", cfgPathAbs, cfgPath)
		}
		return cfgPathAbs, nil
	}

	return cfgPath, nil
}
