package main

import (
    "fmt"
    "net/http"
    "doyoudare-be/internal/config"
    "doyoudare-be/internal/spotify"

    "github.com/gin-gonic/gin"
)

func main() {
    conf, err := config.LoadConfig(".")
    if err != nil {
        fmt.Println("Failed to load config:", err)
        return
    }

    r := gin.Default()

    // Route to start authentication
    r.GET("/login", func(c *gin.Context) {
        authURL, state, err := spotify.StartAuthentication(&conf)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.Redirect(http.StatusTemporaryRedirect, authURL)
		c.SetCookie("oauth_state", state, 3600, "/", "localhost", false, true)
    })

    // Callback route
    r.GET("/callback", func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")
		if code == "" || state == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Code parameter is missing"})
			return
		}
	
		// Retrieve the codeVerifier for this user/session
		// Retrieve the codeVerifier using the state
		codeVerifier, ok := spotify.GetVerifierByState(state)	
		if !ok {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid state parameter"})
            return
        }
		// Clear the saved state and verifier now that it's been used
        spotify.DeleteStateVerifier(state)

		accessToken, refreshToken, err := spotify.ExchangeCodeForToken(&conf, code, codeVerifier)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
			return
		}
	
		// Now, save the accessToken and refreshToken securely
		err = spotify.SaveTokens(accessToken, refreshToken, conf.TokenFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save tokens"})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"message": "Authentication successful"})
	})
	

    r.Run(":8080") // Listen and serve on 0.0.0.0:8080
}
