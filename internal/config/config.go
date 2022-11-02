package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const configFile = "data/config.yaml"

type exchangeRate struct {
	APIKey           string `yaml:"API_key"`
	BaseURI          string `yaml:"base_uri"`
	RefreshRateInMin int64  `yaml:"refresh_rate_in_min"`
}

type tracing struct {
	SamplingRatio float64 `yaml:"sampling_ratio"`
}

type config struct {
	Token        string       `yaml:"token"`
	ExchangeRate exchangeRate `yaml:"exchange_rate"`
	DatabaseDSN  string       `yaml:"db_dsn"`
	ServiceName  string       `yaml:"service_name"`
	Tracing      tracing      `yaml:"tracing"`
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

func (s *Service) ExchangeRateRefreshRateInMin() int64 {
	return s.config.ExchangeRate.RefreshRateInMin
}

func (s *Service) DatabaseDSN() string {
	return s.config.DatabaseDSN
}

func (s *Service) ServiceName() string {
	return s.config.ServiceName
}

func (s *Service) SamplingRatio() float64 {
	return s.config.Tracing.SamplingRatio
}
