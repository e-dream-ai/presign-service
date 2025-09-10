package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	BucketName     string `mapstructure:"bucket_name"`
	APIKey         string `mapstructure:"api_key"`
	Port           string `mapstructure:"port"`
	AWSRegion      string `mapstructure:"aws_region"`
	AWSAccessKeyID string `mapstructure:"aws_access_key_id"`
	AWSSecretKey   string `mapstructure:"aws_secret_access_key"`
	AWSEndpointURL string `mapstructure:"aws_endpoint_url"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("port", "8080")

	err = viper.ReadInConfig()
	if err != nil {
		// It's okay if config file doesn't exist, we can use env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
