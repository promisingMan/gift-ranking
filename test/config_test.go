package test

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"ranking/config"
	"testing"
)

func TestLoadDevConfig(t *testing.T) {
	configuration, err := LoadConfigurationTest("./config-dev.yaml")
	if err != nil {
		fmt.Printf("load configuration error: %v\n", err)
	}
	fmt.Println(configuration)
}

func LoadConfigurationTest(path string) (config *config.Config, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return
	}
	fmt.Println("in func(yaml.v2): ", config)
	return
}
