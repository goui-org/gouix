package config

import (
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port  int    `yaml:"port"`
	Proxy string `yaml:"proxy"`
}

type BuildConfig struct {
	Panic   string `yaml:"panic"`
	Debug   bool   `yaml:"debug"`
	Opt     string `yaml:"opt"`
	WASMOpt bool   `yaml:"wasm_opt"`
	NoTraps bool   `yaml:"no_traps"`
}

type Config struct {
	Server *ServerConfig `yaml:"server"`
	Build  *BuildConfig  `yaml:"build"`
}

func Get() *Config {
	f, err := os.Open("goui.yml")
	if err != nil {
		log.Fatalln(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatalln(err)
	}
	if cfg.Build.Panic == "" {
		cfg.Build.Panic = "trap"
	}
	if cfg.Build.Opt == "" {
		cfg.Build.Opt = "2"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 3000
	}
	return &cfg
}
