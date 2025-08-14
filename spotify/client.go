package spotify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	ID          string
	Secret      string
	RedirectURI string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func (sc *Client) ExchangeCodeForToken(code string) (string, string, int, error) {
	tokenURL := "https://accounts.spotify.com/api/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", sc.RedirectURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(sc.ID + ":" + sc.Secret))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", 0, fmt.Errorf("spotify returned non-200 status: %s", resp.Status)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", "", 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return tokenResponse.AccessToken, tokenResponse.RefreshToken, tokenResponse.ExpiresIn, nil
}

func (sc *Client) RefreshAccessToken(refreshToken string) (string, string, error) {
	tokenURL := "https://accounts.spotify.com/api/token"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(sc.ID + ":" + sc.Secret))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("spotify returned non-200 status: %s", resp.Status)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", "", fmt.Errorf("failed to decode response: %w", err)
	}

	return tokenResponse.AccessToken, tokenResponse.RefreshToken, nil
}
