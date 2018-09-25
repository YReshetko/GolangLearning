package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	defaultConfigFileName = "configuration_dev.json"
)

type Configuration struct {
	HostConfiguration    HostConfig
	EmailStructure       MailStructure
	StorageConfiguration StorageConfig
}
type HostConfig struct {
	ImapHost       string
	ImapPort       int
	ClientEmail    string
	ClientPassword string
}
type MailStructure struct {
	ExpectedSender    string
	FileNameRegExp    string
	WhoCallsRegExp    string
	InputNumberRegExp string
	ParticipantRegExp string
}

type StorageConfig struct {
	DbHost               string
	DbPort               string
	DbName               string
	CollectionName       string
	LocalStorageBasePath string
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
