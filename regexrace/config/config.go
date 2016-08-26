package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// LoadConfig through Viper.
func LoadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("config/")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
}
