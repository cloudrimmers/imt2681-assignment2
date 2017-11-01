package types

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"runtime"
)

// WebConfig configures the web-server from the config.json
type WebConfig struct {
	Currencies        []string
	CollectionWebhook string
	CollectionTick    string
}

// Load the settings file to configure settings
func (v *WebConfig) Load() *WebConfig {

	_, filename, _, _ := runtime.Caller(0)

	data, err := ioutil.ReadFile(path.Dir(filename) + "/webconfig.json")
	if err != nil {
		panic(err.Error())
	}

	if err = json.Unmarshal(data, &v); err != nil {
		panic(err.Error())
	}
	log.Println("Validation settings file: ", v)
	return v
}
