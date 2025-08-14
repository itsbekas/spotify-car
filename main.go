package main

import (
	"os"

	"itsbekas/spotify-car/handlers"
	"itsbekas/spotify-car/spotify"
	"itsbekas/spotify-car/store"
	"itsbekas/spotify-car/util"

	"github.com/gin-gonic/gin"
)

func main() {

	util.CheckEnv()

	sessionStore := store.NewSessionStore()

	spotifyClient := &spotify.Client{
		ID:          os.Getenv("SPOTIFY_CLIENT_ID"),
		Secret:      os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURI: os.Getenv("SPOTIFY_REDIRECT_URI"),
	}

	authHandler := &handlers.AuthHandler{
		Store:   sessionStore,
		Spotify: spotifyClient,
	}

	router := gin.Default()

	router.POST("/register", authHandler.RegisterClient)
	router.GET("/auth", authHandler.StartAuth)
	router.GET("/callback", authHandler.HandleCallback)
	router.GET("/token", authHandler.GetToken)
	router.GET("/refresh", authHandler.RefreshToken)

	router.Run(":8080")
}
