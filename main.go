package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/elliotwutingfeng/the-ghost-of-you-spotify/src"
)

func main() {
	// Load environment variables from .env file
	_, projectRoot, _, _ := runtime.Caller(0)
	pathToEnv := filepath.Join(filepath.Dir(projectRoot), ".env")
	if err := src.LoadDotEnv(pathToEnv); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectHost := os.Getenv("REDIRECT_HOST") // https://developer.mozilla.org/en-US/docs/Web/API/URL/host
	refreshToken := os.Getenv("REFRESH_TOKEN")
	market := os.Getenv("MARKET")

	// Get latest tokens.
	accessToken, refreshToken := src.GetTokens(clientID, clientSecret, redirectHost, refreshToken, pathToEnv)

	// Write latest refresh token to .env file.
	src.UpdateEnvVar(pathToEnv, "REFRESH_TOKEN", refreshToken)

	// Sync liked tracks.
	src.Update(accessToken, market)
}
