package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func getEnvsPath() string {
	_, b, _, _ := runtime.Caller(0)
	configPath := filepath.Dir(b)
	srcPath := filepath.Dir(configPath)
	rootPath := filepath.Dir(srcPath)
	envsPath := filepath.Join(rootPath, "envs")

	return envsPath
}

func LoadEnv() {
	envsPath := getEnvsPath()

	currentFile := "current.env"
	currentEnvPath := filepath.Join(envsPath, currentFile)
	err := godotenv.Load(currentEnvPath)
	if err != nil {
		log.Fatalf("Error loading " + currentFile + " file")
	}

	envName := os.Getenv("ENV")
	envPath := filepath.Join(envsPath, envName)

	firebasePath := filepath.Join(envPath, "firebase.json")
	err = os.Setenv("FIREBASE_PATH", firebasePath)
	if err != nil {
		log.Fatalf("Error setting firebase path file")
	}
}
