package types

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// WebConfig configures the web-server from the config.json
type WebConfig struct {
	Currencies []string
	MinTrigger float64
	MaxTrigger float64
}

// Load the settings file to configure settings
func (v *WebConfig) Load() *WebConfig {
	data, err := ioutil.ReadFile("./webconfig.json")
	if err != nil {
		panic(err.Error())
	}

	if err = json.Unmarshal(data, &v); err != nil {
		panic(err.Error())
	}
	log.Println("Validation settings file: ", v)
	return v
}
