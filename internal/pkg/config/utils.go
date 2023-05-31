package config

import "github.com/spf13/viper"

func NewViperConfig(fileName, configType string) (viper.Viper, error) {
	config := viper.New()
	config.SetConfigType(configType)
	config.AddConfigPath(".")
	config.SetConfigName(fileName)
	if err := config.ReadInConfig(); err != nil {
		return viper.Viper{}, err
	}
	return *config, nil
}
