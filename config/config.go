package config

import (
	"encoding/json"
	"os"
)

// ServerConf is the config struct for the main server properties
type ServerConf struct {
	SSHPort    *string
	SSHKeyPath *string
	SSHKeyName *string

	HTTPPort          *string
	HTTPFileServerDir *string
}

// GameConf is the config struct for the game properties
type GameConf struct {
	Manager *GameManagerConf
	Server  *GameServerConf
	Player  *PlayerConf
}

// Config represents a configuration object for the sshtron game
type Config struct {
	Testkey string

	Server *ServerConf
	Game   *GameConf
}

/*
CreateConfig takes path to a config file as a string,
	and attempts to parse the json into a Config struct
*/
func CreateConfig(cfgFilePath string) (*Config, error) {
	conf := &Config{}

	cfgFile, err := os.Open(cfgFilePath)
	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(cfgFile)
	err = jsonParser.Decode(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
