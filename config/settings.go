package config

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

const settingsPath = "/app/config"

var settingsFiles = map[string]string{
	"development": path.Join(settingsPath, "appConfig.dev.yml"),
	"production":  path.Join(settingsPath, "appConfig.yml"),
}

// AWSSettings contains the settings used for aws
type AWSSettings struct {
	AccessKeyID         string `yaml:"accessKeyID"`
	SecretAccessKey     string `yaml:"secretAccessKey"`
	CognitoUserPoolID   string `yaml:"cognitoUserPoolID"`
	CognitoAppClientID  string `yaml:"cognitoAppClientID"`
	CognitoClientSecret string `yaml:"cognitoClientSecret"`
	Region              string `yaml:"region"`
}

// Settings struct to unmashal config yml setting
type Settings struct {
	Aws AWSSettings
}

// Environment the app is running in either "Production" or "Development"
var Environment string = os.Getenv("Environment")

// LoadSettings loads the settings from the
func LoadSettings() (*Settings, error) {
	config, err := ioutil.ReadFile(settingsFiles[Environment])
	if err != nil {
		return nil, err
	}

	settings := &Settings{}
	err = yaml.Unmarshal(config, settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
