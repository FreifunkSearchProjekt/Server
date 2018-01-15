package config
import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Search struct {
		ModeAppservice bool `yaml:"mode_appservice"`
		Stemming       struct {
		} `yaml:"stemming"`
	} `yaml:"search"`

	Logging struct {
		Level   string `yaml:"level"`
		Logfile string `yaml:"logfile"`
	} `yaml:"logging"`

	Database struct {
		URI string `yaml:"uri"`
	} `yaml:"database"`

	Bleve struct {
		BasePath string `yaml:"base_path"`
	} `yaml:"bleve"`

	Debug struct {
		PProf bool `yaml:"pprof"`
	} `yaml:"debug"`
}

func LoadConfig(path string) (config *Config, err error) {
	var yamlFile []byte
	if yamlFile, err = ioutil.ReadFile(path); err != nil {
		return
	}

	err = yaml.Unmarshal(yamlFile, &config)
	return
}