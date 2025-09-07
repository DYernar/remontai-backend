package config

import (
	"io/ioutil"
	"log"

	"github.com/caarlos0/env/v11"
	_ "github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AppConfig               *AppConfig               `yaml:"app"`
	AppleSigninCredentials  *AppleSigninCredentials  `yaml:"apple_signin_credentials"`
	GoogleSigninCredentials *GoogleSigninCredentials `yaml:"google_signin_credentials"`
	S3Credentials           string                   `env:"S3_CREDENTIALS" envDefault:"default-value-if-not-set"`
}

type AppConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type AppleSigninCredentials struct {
	TeamID     string `yaml:"team_id"`
	ClientID   string `yaml:"client_id"`
	KeyID      string `yaml:"key_id"`
	PrivateKey string `yaml:"private_key"`
}

type GoogleSigninCredentials struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

// Path: config/config.go
// parse config from file
func ParseConfig() (*Config, error) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	var conf Config
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	if err := env.Parse(&conf); err != nil {
		log.Fatalf("%+v", err)
	}

	return &conf, nil
}
