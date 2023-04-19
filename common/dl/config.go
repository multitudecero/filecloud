package dl

import (
	"encoding/json"
	"os"
)

const ConfigFileName = "./app.cfg"

func ReadConfig(cfg interface{}) error {
	data, err := os.ReadFile(ConfigFileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}
	return nil
}
