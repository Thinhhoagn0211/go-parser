package common

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Configuration struct {
	WebserverHost string `mapstructure:"WEBSERVER_HOST"`
	WebserverPort string `mapstructure:"WEBSERVER_PORT"`
	MongoHost     string `mapstructure:"MONGO_HOST"`
	MongoPort     string `mapstructure:"MONGO_PORT"`
}

var (
	Config *Configuration
	Jwt    string
)

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config *Configuration, err error) {
	Config = new(Configuration)
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	Config = config

	return
}
