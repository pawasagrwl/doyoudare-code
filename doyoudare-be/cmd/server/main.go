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
        authURL, _, err := spotify.StartAuthentication(&conf)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.Redirect(http.StatusTemporaryRedirect, authURL)
    })

    // Callback route
    r.GET("/callback", func(c *gin.Context) {
        // This needs to handle the callback and exchange the code
        // Omitted for brevity, assuming implementation is similar to what has been discussed
    })

    r.Run(":8080") // Listen and serve on 0.0.0.0:8080
}
