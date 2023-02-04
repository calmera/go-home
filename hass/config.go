package hass

import (
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
)

func LoadHassConfig(configFile string) ([]*Module, error) {
	cf, err := homedir.Expand(configFile)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(cf)
	if err != nil {
		return nil, err
	}

	return ParseHassModules(b)
}

func ParseHassModules(b []byte) ([]*Module, error) {
	var configs []ModuleConfig
	if err := json.Unmarshal(b, &configs); err != nil {
		return nil, err
	}

	modules := make([]*Module, len(configs))
	for idx, config := range configs {
		m, err := NewModule(config)
		if err != nil {
			return nil, err
		}

		modules[idx] = m
	}

	return modules, nil
}

type ModuleConfig struct {
	Id    string `json:"id"`
	Url   string `json:"url"`
	Token string `json:"token"`
}
