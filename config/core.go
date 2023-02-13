// Package config houses configurations for application.
package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// Configuration will allow the application to run with different settings.
type Configuration struct {
	Port            string `json:"port"`
	DefaultLanguage string `json:"default_language"`
	LegacyEndpoint  string `json:"legacy_endpoint"`
	DatabaseType    string `json:"database_type"`
	DatbaseURL      string `json:"database_url"`
}

var defaultConfiguration = Configuration{
	Port:            ":8080",
	DefaultLanguage: "english",
}

// LoadFromEnv will load configuration solely from environment.
func (c *Configuration) LoadFromEnv() {
	if lang := os.Getenv("DEFAULT_LANGUAGE"); lang != "" {
		c.DefaultLanguage = lang
	}
	if port := os.Getenv("PORT"); port != "" {
		c.Port = port
	}
}

// ParsePort will check to see if the port is  in the proper format and a number.
func (c *Configuration) ParsePort() {
	if c.Port[0] != ':' {
		c.Port = ":" + c.Port
	}
	if _, err := strconv.Atoi(string(c.Port[1:])); err != nil {
		fmt.Printf("invalid port %s", c.Port)
		c.Port = defaultConfiguration.Port
	}
}

// LoadFromJSON will read a JSON file and update configuration based on the file.
func (c *Configuration) LoadFromJSON(path string) error {
	log.Printf("loading configuration from file: %s\n", path)
	b, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		log.Printf("unable to load file: %s\n", err.Error())
		return errors.New("unable to load configuration")
	}
	if err := json.Unmarshal(b, c); err != nil {
		log.Printf("unable to parse file: %s\n", err.Error())
		return errors.New("unable to load configuration")
	}
	// Verify required fields.
	if c.Port == "" {
		log.Printf("empty port, reverting to default")
		c.Port = defaultConfiguration.Port
	}
	if c.DefaultLanguage == "" {
		log.Printf("empty language, reverting to default")
		c.DefaultLanguage = defaultConfiguration.DefaultLanguage
	}
	return nil
}

// LoadConfiguration will provide cycle through flags, files
// and finally environment variables to load configuration.
func LoadConfiguration() Configuration {
	cfgFileFlag := flag.String("config_file", "", "load configurations from a file")
	portFlag := flag.String("port", "", "set port")

	flag.Parse()
	cfg := defaultConfiguration

	if cfgFileFlag != nil && *cfgFileFlag != "" {
		if err := cfg.LoadFromJSON(*cfgFileFlag); err != nil {
			log.Printf("unable to load configuration from json: %s, using default values\n", *cfgFileFlag)
		}
	}

	cfg.LoadFromEnv()

	if portFlag != nil && *portFlag != "" {
		cfg.Port = *portFlag
	}

	cfg.ParsePort()
	return cfg
}
