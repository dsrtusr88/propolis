package main

import (
	"io/ioutil"
	"sync"

	"github.com/pkg/errors"
	"gitlab.com/catastrophic/assistance/logthis"
	yaml "gopkg.in/yaml.v2"
)

var config *Config
var onceConfig sync.Once

const (
	DefaultConfigurationFile = "config.yaml"
	errorReadingConfig       = "could not read configuration file"
	errorLoadingYAML         = "could not decode yaml"
)

type Config struct {
	General       *ConfigGeneral
	IRC           *ConfigIrc
	IrcConfigured bool
}

func NewConfig(path string) (*Config, error) {
	var newConfigErr error
	onceConfig.Do(func() {
		// TODO check path has yamlExt!
		newConf := &Config{}
		if err := newConf.Load(path); err != nil {
			newConfigErr = err
			return
		}
		// set the global pointer once everything is OK.
		config = newConf
	})
	return config, newConfigErr
}

func (c *Config) String() string {
	txt := c.General.String() + "\n"
	if c.IrcConfigured {
		txt += c.IRC.String() + "\n"
	}
	return txt
}

func (c *Config) Load(file string) error {
	// loading the configuration file
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, errorReadingConfig)
	}
	return c.LoadFromBytes(b)
}

func (c *Config) LoadFromBytes(b []byte) error {
	err := yaml.Unmarshal(b, &c)
	if err != nil {
		return errors.Wrap(err, errorLoadingYAML)
	}
	return c.check()
}

func (c *Config) check() error {
	// general checks
	if c.General == nil {
		return errors.New("General configuration required")
	}
	if err := c.General.check(); err != nil {
		return errors.Wrap(err, "Error reading general configuration")
	}
	// setting log level
	logthis.SetLevel(c.General.LogLevel)
	logthis.SetStdOutput(true)

	if c.IRC != nil {
		if err := c.IRC.check(); err != nil {
			return errors.Wrap(err, "Error reading Metadata configuration")
		}
		c.IrcConfigured = true
	}
	return nil
}
