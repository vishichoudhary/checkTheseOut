package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

type SystemSettings struct {
	Configurations `json:"configurations"`
}

func GetSystemSettings() SystemSettings {
	return *systemSettings
}

func (conf Configurations) SetLogLevel() log.Level {
	switch conf.Logging.Level {
	case "DEBUG":
		return log.DebugLevel
	default:
		return log.InfoLevel
	}
}

type Service struct {
	APIVersion string `json:"api_version"`
	Port       int    `json:"port"`
}
type Functional struct {
	InactivityTimeInSec int `json:"inactivityTimeInSec"`
}
type Logging struct {
	Level string `json:"level"`
}
type Configurations struct {
	Service    Service    `json:"service"`
	Functional Functional `json:"functional"`
	Logging    Logging    `json:"logging"`
}

var once sync.Once
var systemSettings *SystemSettings

func Load(file string) *SystemSettings {
	once.Do(func() {
		jsonFile, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		json.Unmarshal(byteValue, &systemSettings)
	})
	return systemSettings
}
