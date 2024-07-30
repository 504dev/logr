package config

import (
	"crypto/sha256"
	"encoding/base64"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

type Args struct {
	Configpath string
	Retries    int
}

type ConfigData struct {
	Bind struct {
		Http string `yaml:"http"`
		Udp  string `yaml:"udp"`
		Grpc string `yaml:"grpc"`
		Prom string `yaml:"prom"`
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
func (c *ConfigData) NeedSetup() bool {
	return c.OAuth.Github.ClientId == "" || c.OAuth.Github.ClientSecret == ""
}

type Config struct {
	sync.RWMutex
	Data ConfigData
}

func (c *Config) FromFile(configpath string) error {
	yamlFile, err := os.ReadFile(configpath)
	if err != nil {
		return err
	}
	yamlFile = []byte(os.ExpandEnv(string(yamlFile)))
	err = yaml.Unmarshal(yamlFile, &c.Data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Get() *ConfigData {
	c.RLock()
	defer c.RUnlock()
	return &c.Data
}

func (c *Config) Set(set func(c *ConfigData)) {
	c.Lock()
	clone := c.Data
	set(&clone)
	c.Data = clone
	c.Unlock()
}

func (c *Config) ToFile(configpath string) error {
	d, err := yaml.Marshal(&c.Data)
	if err != nil {
		return err
	}
	err = os.WriteFile(configpath, d, 0644)
	if err != nil {
		return err
	}
	return nil
}
