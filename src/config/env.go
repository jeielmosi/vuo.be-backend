package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/joho/godotenv"
)

const CURRENT_ENV = "CURRENT_ENV"

func getEnvsPath() string {
	_, b, _, _ := runtime.Caller(0)
	configPath := filepath.Dir(b)
	srcPath := filepath.Dir(configPath)
	rootPath := filepath.Dir(srcPath)
	envsPath := filepath.Join(rootPath, "envs")

	return envsPath
}

func LoadEnv(envName string) {

	envsPath := getEnvsPath()
	envPath := filepath.Join(envsPath, envName)

	firebasePath := filepath.Join(envPath, "firebase.json")
	err := os.Setenv(envName+"_FIREBASE_PATH", firebasePath)
	if err != nil {
		log.Fatalf("Error setting firebase path file")
		os.Exit(1)
	}

	dotEnvPath := filepath.Join(envPath, ".env")
	err = godotenv.Overload(dotEnvPath)
	if err != nil {
		log.Fatalf("Error loading " + dotEnvPath + " file")
		os.Exit(1)
	}
}

func GetEnv() string {
	return os.Getenv(CURRENT_ENV)
}

var once sync.Once

func loadDefaulfEnv() {
	currentFile := "current.env"
	var err error = nil
	once.Do(
		func() {
			envsPath := getEnvsPath()
			currentEnvPath := filepath.Join(envsPath, currentFile)
			err = godotenv.Load(currentEnvPath)
		},
	)

	if err != nil {
		log.Fatalf("Error loading " + currentFile + " file")
		os.Exit(1)
	}
}

func Load() {
	loadDefaulfEnv()
	envName := GetEnv()
	LoadEnv(envName)
}
