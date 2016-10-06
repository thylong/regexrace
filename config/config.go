package config

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// LoadConfig Load config.yml content with Viper and override with env variables.
func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("config/")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	for _, settingKey := range viper.AllKeys() {
		viper.BindEnv(settingKey, strings.ToUpper(settingKey))
	}
}
