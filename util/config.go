package util

import (
	"os"

	"github.com/saurlax/net-vigil/tix"
	"gopkg.in/yaml.v3"
)

type NetVigilConfig struct {
	CaptureInterval int            `yaml:"capture_interval"`
	CheckInterval   int            `yaml:"check_interval"`
	Buffer          int            `yaml:"buffer"`
	Port            int            `yaml:"port"`
	Local           tix.Local      `yaml:"local"`
	ThreatBook      tix.ThreatBook `yaml:"threatbook"`
}

var Config NetVigilConfig

func init() {
	data, _ := os.ReadFile("config.yml")
	yaml.Unmarshal(data, &Config)
}
