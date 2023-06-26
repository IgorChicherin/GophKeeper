package config

import log "github.com/sirupsen/logrus"

type ClientConfig struct {
	Address string `mapstructure:"address"`
}

func GetClientConfig() (ClientConfig, error) {
	v, err := NewViperConfig("client", "json")

	if err != nil {
		return ClientConfig{}, err
	}

	cnf := ClientConfig{}
	if err := v.Unmarshal(&cnf); err != nil {
		log.Errorln(err)
		return ClientConfig{}, err
	}
	return cnf, nil
}
