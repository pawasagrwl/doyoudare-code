package spotify

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"encoding/json"
    "io"

	"doyoudare-be/internal/config"
)

// Generate a secure, random string for the PKCE code verifier
func generateCodeVerifier() (string, error) {
	randomBytes := make([]byte, 32) // 32 bytes give us 256 bits of entropy
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes), nil
}

// Generate the code challenge based on the code verifier
func generateCodeChallenge(verifier string) (string, error) {
	// S256 method
	h := sha256.New()
	if _, err := h.Write([]byte(verifier)); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(h.Sum(nil)), nil
}

// Construct the Spotify authorization URL with PKCE
func buildAuthURL(conf *config.Config, codeChallenge string) string {
	// Define the required query parameters
	params := url.Values{}
	params.Add("client_id", conf.SpotifyClientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", conf.SpotifyRedirectURI)
	params.Add("code_challenge_method", "S256")
	params.Add("code_challenge", codeChallenge)
	params.Add("scope", "user-read-private user-read-currently-playing user-read-playback-state playlist-read-private")

	// Construct the full URL
	authURL := fmt.Sprintf("https://accounts.spotify.com/authorize?%s", params.Encode())
	return authURL
}

// StartAuthentication initiates the authentication process
func StartAuthentication(conf *config.Config) (authURL string, codeVerifier string, err error) {
	codeVerifier, err = generateCodeVerifier()
	if err != nil {
		return "", "", err
	}
	codeChallenge, err := generateCodeChallenge(codeVerifier)
	if err != nil {
		return "", "", err
	}
	authURL = buildAuthURL(conf, codeChallenge)
	return authURL, codeVerifier, nil
}

type TokenResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token"`
    Scope        string `json:"scope"`
}

func parseTokenResponse(body io.Reader) (accessToken string, refreshToken string, err error) {
    var resp TokenResponse
    err = json.NewDecoder(body).Decode(&resp)
    if err != nil {
        return "", "", err
    }
    return resp.AccessToken, resp.RefreshToken, nil
}

// Exchange the authorization code for access and refresh tokens
func ExchangeCodeForToken(conf *config.Config, code string, codeVerifier string) (accessToken string, refreshToken string, err error) {
	// Define the request body
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", conf.SpotifyRedirectURI)
	data.Set("client_id", conf.SpotifyClientID)
	data.Set("code_verifier", codeVerifier)

	// Make the request
	resp, err := http.PostForm("https://accounts.spotify.com/api/token", data)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Parse the response (omitted for brevity)
	// Assume we have a function parseTokenResponse that extracts the tokens from the response

	return parseTokenResponse(resp.Body)
}




