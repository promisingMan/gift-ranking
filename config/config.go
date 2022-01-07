package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var AppConfig *Config

type Config struct {
	Redis   RedisConfig   `yaml:"redis"`
	MongoDb MongoDbConfig `yaml:"mongodb"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

type MongoDbConfig struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
}

// LoadConfiguration 加载配置文件，在main入口处执行一次即可
func LoadConfiguration(env string) error {
	if env == "" {
		env = "dev"
	}
	path := fmt.Sprintf("./config/config-%s.yaml", env)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &AppConfig)
	if err != nil {
		return err
	}
	log.Printf("load configuration successfully : %v", AppConfig)
	return nil
}
