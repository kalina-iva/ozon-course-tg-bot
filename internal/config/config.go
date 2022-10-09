package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const configFile = "data/config.yaml"

type exchangeRate struct {
	APIKey       string `yaml:"API_key"`
	BaseURI      string `yaml:"base_uri"`
	TimeoutInMin int64  `yaml:"timeout_in_min"`
}

type config struct {
	Token        string       `yaml:"token"`
	ExchangeRate exchangeRate `yaml:"exchange_rate"`
}

type Service struct {
	config config
}

func New() (*Service, error) {
	s := &Service{}

	rawYAML, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &s.config)
	if err != nil {
		return nil, errors.Wrap(err, "parsing yaml")
	}

	return s, nil
}

func (s *Service) Token() string {
	return s.config.Token
}

func (s *Service) ExchangeRateAPIKey() string {
	return s.config.ExchangeRate.APIKey
}

func (s *Service) ExchangeRateBaseURI() string {
	return s.config.ExchangeRate.BaseURI
}

func (s *Service) ExchangeRateTimeout() int64 {
	return s.config.ExchangeRate.TimeoutInMin
}
