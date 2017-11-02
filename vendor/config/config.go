package config

import (
	"fmt"
	"path"
	"runtime"

	"github.com/BurntSushi/toml"
	"log"
)

var (
	config Config
)

type Config struct {
	DbDsn string
	Port  string
}

func init() {
	_, confFile, _, _ := runtime.Caller(1)

	_, err := toml.DecodeFile(path.Dir(confFile)+"/config.toml", &config)

	if err != nil {
		log.Fatal(fmt.Printf("config file error: %s", err))
	}
}

func GetConfig() *Config {
	return &config
}
