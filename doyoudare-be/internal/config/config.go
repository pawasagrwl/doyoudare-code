package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    SpotifyClientID     string `mapstructure:"SPOTIFY_CLIENT_ID"`
    SpotifyClientSecret string `mapstructure:"SPOTIFY_CLIENT_SECRET"`
    SpotifyRedirectURI  string `mapstructure:"SPOTIFY_REDIRECT_URI"`
    TokenFilePath       string `mapstructure:"TOKEN_FILE_PATH"`
}

func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("env")
    viper.SetConfigType("env")

    viper.AutomaticEnv()

    err = viper.ReadInConfig()
    if err != nil {
        return
    }

    err = viper.Unmarshal(&config)
    return
}
