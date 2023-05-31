package config

import (
	log "github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Database string `json:"database"`
	Address  string `json:"address"`
	Key      string `json:"key"`
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
