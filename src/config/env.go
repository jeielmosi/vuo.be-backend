package config

import (
	"path/filepath"
	"runtime"
)

func getEnvsPath() string {
	_, b, _, _ := runtime.Caller(0)
	configPath := filepath.Dir(b)
	srcPath := filepath.Dir(configPath)
	rootPath := filepath.Dir(srcPath)
	envsPath := filepath.Join(rootPath, "envs")

	return envsPath
}
