package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"itsbekas/spotify-car/spotify"
	"itsbekas/spotify-car/store"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Store   *store.SessionStore
	Spotify *spotify.Client
}

func (h *AuthHandler) RegisterClient(c *gin.Context) {
	picoID := c.Query("pico_id")
	sessionHash := c.Query("session_hash")

	if picoID == "" || sessionHash == "" {
		c.String(http.StatusBadRequest, "Missing pico_id or session_hash")
		return
	}

	h.Store.SetInitialData(picoID, sessionHash)
}

func (h *AuthHandler) StartAuth(c *gin.Context) {
	picoID := c.Query("pico_id")

	if picoID == "" {
		c.String(http.StatusBadRequest, "Missing pico_id")
		return
	}

	spotifyAuthURL := "https://accounts.spotify.com/authorize"
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", h.Spotify.ID)
	params.Add("scope", "user-read-currently-playing user-modify-playback-state")
	params.Add("redirect_uri", h.Spotify.RedirectURI)
	params.Add("state", picoID)

	c.Redirect(http.StatusFound, spotifyAuthURL+"?"+params.Encode())
}

func (h *AuthHandler) HandleCallback(c *gin.Context) {
	code := c.Query("code")
	picoID := c.Query("state")

	if code == "" || picoID == "" {
		c.String(http.StatusBadRequest, "Invalid callback request")
		return
	}

	accessToken, refreshToken, expiresIn, err := h.Spotify.ExchangeCodeForToken(code)

	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error exchanging code: %v", err))
		return
	}

	h.Store.SetTokens(picoID, accessToken, refreshToken, expiresIn)

	c.String(http.StatusOK, "Success! You can now press play on your Pico.")
}

func (h *AuthHandler) GetToken(c *gin.Context) {
	picoID := c.Query("pico_id")
	sessionHash := c.Query("session_hash")

	accessToken, refreshToken, found := h.Store.GetTokens(picoID, sessionHash)
	if !found {
		c.String(http.StatusNotFound, "Tokens not ready, invalid or expired.")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
