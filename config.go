package exporter

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type LogConfiguration struct {
	Format string `yaml:"format,omitempty"`
	Level  string `yaml:"level,omitempty"`
}

type HistoricalRatedUsage struct {
	Enable              bool `yaml:"enable,omitempty"`
	NumberOfPriorMonths int  `yaml:"numberOfPriorMonths,omitempty"`
}

type BmcApiConfiguration struct {
	ClientId             string               `yaml:"clientId"`
	ClientSecret         string               `yaml:"clientSecret"`
	TokenUrl             string               `yaml:"tokenUrl"`
	Log                  LogConfiguration     `yaml:"log,omitempty"`
	HistoricalRatedUsage HistoricalRatedUsage `yaml:"historicalRatedUsage,omitempty"`
}

func (m *BmcApiConfiguration) SetLogFormat(logFormat string) *BmcApiConfiguration {
	value := setConfigValueString(m.Log.Format, logFormat, "LOG_FORMAT", "text")

	switch strings.ToLower(value) {
	case "json":
		fallthrough
	case "text":
		m.Log.Format = strings.ToLower(value)
	default:
		m.Log.Format = "text"
	}

	return m
}

func (m *BmcApiConfiguration) SetLogLevel(logLevel string) *BmcApiConfiguration {
	value := setConfigValueString(m.Log.Level, logLevel, "LOG_LEVEL", "info")

	switch strings.ToLower(value) {
	case "panic":
		fallthrough
	case "fatal":
		fallthrough
	case "warning":
		fallthrough
	case "info":
		fallthrough
	case "debug":
		fallthrough
	case "trace":
		m.Log.Level = strings.ToLower(value)
	default:
		m.Log.Level = "info"
	}

	return m
}

func (m BmcApiConfiguration) GetLogFormatter() log.Formatter {
	switch m.Log.Format {
	case "json":
		return &log.JSONFormatter{}
	default:
		return &log.TextFormatter{}
	}
}

func (m BmcApiConfiguration) GetLogLevel() log.Level {
	level, _ := log.ParseLevel(m.Log.Level)
	return level
}

func (m BmcApiConfiguration) ToClientCredentials() clientcredentials.Config {
	config := clientcredentials.Config{
		ClientID:     m.ClientId,
		ClientSecret: m.ClientSecret,
		TokenURL:     m.TokenUrl,
		Scopes:       []string{"bmc", "bmc.read"},
	}

	return config
}

func (m BmcApiConfiguration) String() string {
	return fmt.Sprintf("{ ClientId: %s, ClientSecret: %s, TokenUrl: %s }", m.ClientId, m.ClientSecret, m.TokenUrl)
}

func (m *BmcApiConfiguration) load(configPath string, clientId string, clientSecret string, tokenUrl string) {
	if configPath != "" {
		if err := m.loadFromFile(configPath); err != nil {
			log.Error("Unable to read config file.")
			return
		}
	} else {
		log.WithField("configFile", configPath).Info("Attempted to load PhoenixNAP BMC API configuration from file, but no file was specified.")
	}

	m.ClientId = setConfigValueString(m.ClientId, clientId, "PNAP_BMC_API_CLIENT_ID", "")
	m.ClientSecret = setConfigValueString(m.ClientSecret, clientSecret, "PNAP_BMC_API_CLIENT_SECRET", "")
	m.TokenUrl = setConfigValueString(m.TokenUrl, tokenUrl, "PNAP_BMC_API_TOKEN_URL", "https://auth.phoenixnap.com/auth/realms/BMC/protocol/openid-connect/token")
}

func (m *BmcApiConfiguration) loadFromFile(configPath string) error {
	data, err := os.ReadFile(configPath)

	if err != nil {
		log.WithField("configFile", configPath).WithError(err).Error("Unable to read BMC API Configuration file.")
		return err
	}

	if err := yaml.Unmarshal(data, &m); err != nil {
		log.WithField("configFile", configPath).WithError(err).Error("Unable to unmarshal BMC API Configuration file.")
		return err
	}

	return nil
}

func setConfigValueString(fromFile string, fromArg string, envVarName string, defaultValue string) string {
	fromEnv := os.Getenv(envVarName)

	if fromArg != "" {
		return fromArg
	}

	if fromEnv != "" {
		return fromEnv
	}

	if fromFile != "" {
		return fromFile
	}

	return defaultValue
}

func NewBmcApiConfiguration(configPath string, clientId string, clientSecret string, tokenUrl string) *BmcApiConfiguration {
	var m BmcApiConfiguration
	m.load(configPath, clientId, clientSecret, tokenUrl)
	return &m
}
