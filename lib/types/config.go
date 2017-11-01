package types

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// WebConfig configures the web-server from the config.json
type WebConfig struct {
	Currencies        []string
	CollectionWebhook string
	CollectionFixer   string
}

// Load the settings file to configure settings
func (v *WebConfig) Load() *WebConfig {
	basepath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	data, err := ioutil.ReadFile(basepath + "/webconfig.json")
	log.Println("loading config from : ", basepath+"/webconfig.json")

	if err != nil {
		panic(err.Error())
	}

	if err = json.Unmarshal(data, &v); err != nil {
		panic(err.Error())
	}
	//	log.Println("Validation settings file: ", v)
	return v
}
