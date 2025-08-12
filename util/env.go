package util

import (
	"log"
	"os"
)

func checkVar(varName string) bool {
	val := os.Getenv(varName)
	if val == "" {
		log.Fatalf("%v is not set.", varName)
		return false
	}
	return true
}

func CheckEnv() {
	vars := [3]string{"SPOTIFY_CLIENT_ID", "SPOTIFY_CLIENT_SECRET", "SPOTIFY_REDIRECT_URI"}

	varsExist := true

	for _, v := range vars {
		varsExist = varsExist && checkVar(v)
	}

	if !varsExist {
		os.Exit(1)
	}
}
