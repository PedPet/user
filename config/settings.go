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

// DBSettings contains the settings used for the database connection
type DBSettings struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}

// Validation contains the centrealized settings for validation
type Validation struct {
	Password Password `json:"password"`
}

// Password contains the validation config
type Password struct {
	Required bool `json:"required"`
	Length   struct {
		Min int `json:"min"`
		Max int `json:"max"`
	} `json:"length"`
	Regex []string `json:"regex"`
}

// Settings struct to unmarshal config yml setting
type Settings struct {
	Aws AWSSettings
	DB  DBSettings
}

var environment string = os.Getenv("Environment")

// LoadSettings loads the settings from the yml file
func LoadSettings() (*Settings, error) {
	config, err := ioutil.ReadFile(settingsFiles[environment])
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
