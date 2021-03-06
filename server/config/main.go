/*

Config is a simple library that manages config set-up for servers based on a config file.

*/
package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Dev  *ConfigMode
	Prod *ConfigMode
}

type ConfigMode struct {
	AllowedOrigins    string
	DefaultPort       string
	FirebaseProjectId string
	AdminUserIds      []string
	//This is a dangerous config. Only enable in Dev!
	DisableAdminChecking bool
	StorageConfig        map[string]string
}

func (c *Config) validate() error {

	if c.Dev == nil && c.Prod == nil {
		return errors.New("Neither dev nor prod configuration was valid")
	}

	if c.Dev != nil {
		if err := c.Dev.validate(true); err != nil {
			return err
		}
	}
	if c.Prod != nil {
		return c.Prod.validate(false)
	}
	return nil
}

func (c *ConfigMode) validate(isDev bool) error {
	if c.DefaultPort == "" {
		return errors.New("No default port provided")
	}
	//AllowedOrigins will just be default allow
	if c.AllowedOrigins == "" {
		log.Println("No AllowedOrigins found. Defaulting to '*'")
		c.AllowedOrigins = "*"
	}
	if c.StorageConfig == nil {
		c.StorageConfig = make(map[string]string)
	}
	if c.DisableAdminChecking && !isDev {
		return errors.New("DisableAdminChecking enabled in prod, which is illegal")
	}
	return nil
}

func (c *ConfigMode) OriginAllowed(origin string) bool {

	originUrl, err := url.Parse(origin)

	if err != nil {
		return false
	}

	if c.AllowedOrigins == "" {
		return false
	}
	if c.AllowedOrigins == "*" {
		return true
	}
	allowedOrigins := strings.Split(c.AllowedOrigins, ",")
	for _, allowedOrigin := range allowedOrigins {
		u, err := url.Parse(allowedOrigin)

		if err != nil {
			continue
		}

		if u.Scheme == originUrl.Scheme && u.Host == originUrl.Host {
			return true
		}
	}
	return false
}

const (
	configFileName       = "config.SECRET.json"
	sampleConfigFileName = "config.SAMPLE.json"
)

func Get() (*Config, error) {

	fileNameToUse := configFileName

	if _, err := os.Stat(configFileName); os.IsNotExist(err) {

		if _, err := os.Stat(sampleConfigFileName); os.IsNotExist(err) {
			return nil, errors.New("Couldn't find a " + configFileName + " in current directory (or a SAMPLE). This file is required. Copy a starter one from boardgame/server/api/config.SAMPLE.json")
		}

		fileNameToUse = sampleConfigFileName

	}

	contents, err := ioutil.ReadFile(fileNameToUse)

	if err != nil {
		return nil, errors.New("Couldn't read config file: " + err.Error())
	}

	var config Config

	if err := json.Unmarshal(contents, &config); err != nil {
		return nil, errors.New("couldn't unmarshal config file: " + err.Error())
	}

	if err := config.validate(); err != nil {
		return nil, errors.New("Couldn't validate config: " + err.Error())
	}

	return &config, nil

}
