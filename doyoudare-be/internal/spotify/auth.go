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
    "os"

	"doyoudare-be/internal/config"
)

type Tokens struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

// SaveTokens writes the tokens to a specified file path securely
func SaveTokens(accessToken, refreshToken, filePath string) error {
    tokens := Tokens{AccessToken: accessToken, RefreshToken: refreshToken}
    file, err := json.MarshalIndent(tokens, "", " ")
    if err != nil {
        return err
    }

    return os.WriteFile(filePath, file, 0644)
}

// RefreshAccessToken uses the refresh token to get a new access token
func RefreshAccessToken(refreshToken string, conf *config.Config) (newAccessToken string, err error) {
    data := url.Values{}
    data.Set("grant_type", "refresh_token")
    data.Set("refresh_token", refreshToken)
    data.Set("client_id", conf.SpotifyClientID)

    resp, err := http.PostForm("https://accounts.spotify.com/api/token", data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    accessToken, _, err := parseTokenResponse(resp.Body) // Refresh token response may not include a new refresh token
    return accessToken, err
}

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
func buildAuthURL(conf *config.Config, codeChallenge, state string) string {	// Define the required query parameters
	params := url.Values{}
	params.Add("client_id", conf.SpotifyClientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", conf.SpotifyRedirectURI)
	params.Add("code_challenge_method", "S256")
	params.Add("code_challenge", codeChallenge)
	params.Add("scope", "user-read-private user-read-currently-playing user-read-playback-state playlist-read-private")
	params.Add("state", state)
	// Construct the full URL
	authURL := fmt.Sprintf("https://accounts.spotify.com/authorize?%s", params.Encode())
	return authURL
}

// StartAuthentication initiates the authentication process
// StartAuthentication initiates the authentication process and returns the URL and state
func StartAuthentication(conf *config.Config) (authURL, state string, err error) {
	state, err = generateState()
	if err != nil {
		return "", "", err
	}

	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return "", "", err
	}
	codeChallenge, err := generateCodeChallenge(codeVerifier)
	if err != nil {
		return "", "", err
	}

	// Store the codeVerifier with the state
	SaveStateVerifier(state, codeVerifier)

	authURL = buildAuthURL(conf, codeChallenge, state)
	return authURL, state, nil
}

// Add a helper function to generate a state parameter
func generateState() (string, error) {
	randomBytes := make([]byte, 16) // 128 bits is a common length for a state parameter
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes), nil
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




