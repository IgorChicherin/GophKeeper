package config

import (
	log "github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Database        string `mapstructure:"database"`
	Address         string `mapstructure:"address"`
	Key             string `mapstructure:"key"`
	PrivateCertPath string `mapstructure:"private_cert"`
	PublicCertPath  string `mapstructure:"public_cert"`
}

func GetServerConfig() (ServerConfig, error) {
	v, err := NewViperConfig("server", "json")

	if err != nil {
		return ServerConfig{}, err
	}

	cnf := ServerConfig{}
	if err := v.Unmarshal(&cnf); err != nil {
		log.Errorln(err)
		return ServerConfig{}, err
	}
	return cnf, nil
}
