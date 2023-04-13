package main

import (
	"fmt"

	"github.com/spf13/viper"

	"MovieTinder/internal/trakt"
)

// BuildSettings reads the Settings from the config file.
func BuildSettings() *trakt.Settings {
	viper.SetConfigName("stg_settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return &trakt.Settings{
		ClientID:     viper.GetString("trakt.client_id"),
		ClientSecret: viper.GetString("trakt.client_secret"),
		APIVersion:   viper.GetString("trakt.api_version"),
		Domain:       viper.GetString("trakt.domain"),
	}
}
