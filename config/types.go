package config

import (
	"crypto/sha256"
	"encoding/base64"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

type CommandLineArgs struct {
	Configpath string
	Retries    int
}

type ConfigData struct {
	Bind struct {
		HTTP string `yaml:"http"`
		UDP  string `yaml:"udp"`
		GRPC string `yaml:"grpc"`
	} `yaml:"bind"`
	OAuth struct {
		Github struct {
			Org          string `yaml:"org"           json:"org"`
			ClientId     string `yaml:"client_id"     json:"client_id"`
			ClientSecret string `yaml:"client_secret" json:"-"`
		} `yaml:"github"`
		JwtSecret string `yaml:"jwt_secret"`
	} `yaml:"oauth"`
	Clickhouse    string `yaml:"clickhouse"`
	Mysql         string `yaml:"mysql"`
	Redis         string `yaml:"redis"`
	AllowNoCipher bool   `yaml:"allow_no_cipher"`
	DemoDash      struct {
		Enabled bool   `yaml:"enabled"`
		Llm     string `yaml:"llm"`
	} `yaml:"demo_dash"`
	RecaptchaSecret string `yaml:"recaptcha"`
}

func (c *ConfigData) GetJwtSecret() string {
	if c.OAuth.JwtSecret != "" {
		return c.OAuth.JwtSecret
	}
	hash := sha256.Sum256([]byte(c.OAuth.Github.ClientSecret))
	return base64.StdEncoding.EncodeToString(hash[:])
}
func (c *ConfigData) IsSetupRequired() bool {
	return c.OAuth.Github.ClientId == "" || c.OAuth.Github.ClientSecret == ""
}

type Config struct {
	mu         sync.RWMutex
	data       ConfigData
	configpath string
}

func (c *Config) ReadFromFile(configpath string) error {
	c.configpath = configpath
	yamlFile, err := os.ReadFile(configpath)
	if err != nil {
		return err
	}
	yamlFile = []byte(os.ExpandEnv(string(yamlFile)))
	err = yaml.Unmarshal(yamlFile, &c.data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Get() *ConfigData {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return &c.data
}

func (c *Config) Set(setterFunc func(c *ConfigData)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	clone := c.data
	setterFunc(&clone)
	c.data = clone
}

func (c *Config) Save() error {
	return c.SaveToFile(c.configpath)
}

func (c *Config) SaveToFile(configpath string) error {
	data, err := yaml.Marshal(&c.data)
	if err != nil {
		return err
	}
	return os.WriteFile(configpath, data, 0644)
}
