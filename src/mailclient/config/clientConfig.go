package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	defaultConfigFileName = "configuration.json"
)

type Configuration struct {
	ImapHost       string
	ImapPort       int
	ClientEmail    string
	ClientPassword string
}

type ErrorLoadConfig struct {
	fileName string
}

func (err ErrorLoadConfig) Error() string {
	return fmt.Sprintf("Couldn't load file %s", err.fileName)
}

func (config *Configuration) Load() error {
	return config.LoadWithFileName(defaultConfigFileName)
}
func (config *Configuration) LoadWithFileName(fileName string) error {
	file, ok := ioutil.ReadFile(fileName)
	if ok != nil {
		err := ErrorLoadConfig{fileName}
		return err
	}
	json.Unmarshal(file, config)
	return nil
}