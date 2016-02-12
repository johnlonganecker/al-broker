package config

import (
	"io/ioutil"

	"net/http"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port   string  `yaml:"port"`
	SbPort string  `yaml:"sb_port"`
	SbUrl  string  `yaml:"sb_url"`
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Listen      Transaction `yaml:"listen"`
	Destination Transaction `yaml:"destination"`
}

type Transaction struct {
	Url         string            `yaml:"url,omitempty"`
	HttpMethod  string            `yaml:"http_method"`
	Headers     http.Header       `yaml:"headers"`
	Mappings    map[string]string `yaml:"mappings"`
	ExtraFields map[string]string `yaml:"extra_fields"`
}

func (c *Config) LoadConfigFile(filepath string) error {

	var data []byte

	data, err := LoadFile(filepath)
	if err != nil {
		return err
	}

	if err := Unmarshal(c, data); err != nil {
		return err
	}

	return nil
}

func LoadFile(filepath string) ([]byte, error) {

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func Unmarshal(c *Config, data []byte) error {

	// unmarshal yaml
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}

	return nil
}
