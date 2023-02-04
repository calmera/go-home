package hass

import (
	"encoding/json"
	"github.com/mitchellh/go-homedir"
	"github.com/nats-io/nats.go"
	"io/ioutil"
)

func LoadHassConfig(configFile string) ([]*HassModule, error) {
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

func ParseHassModules(b []byte) ([]*HassModule, error) {
	var configs []HassModuleConfig
	if err := json.Unmarshal(b, &configs); err != nil {
		return nil, err
	}

	modules := make([]*HassModule, len(configs))
	for idx, config := range configs {
		m, err := NewModule(config)
		if err != nil {
			return nil, err
		}

		modules[idx] = m
	}

	return modules, nil
}

func NewModule(config HassModuleConfig) (*HassModule, error) {
	return &HassModule{config: config}, nil
}

type HassModuleConfig struct {
	Id    string `json:"id"`
	Url   string `json:"url"`
	Token string `json:"token"`
}

type HassModule struct {
	config HassModuleConfig
}

func (m *HassModule) Identifier() string {
	return m.config.Id
}

func (m *HassModule) Run(nc *nats.Conn) error {
	return nil
}

func (m *HassModule) Close() {

}
