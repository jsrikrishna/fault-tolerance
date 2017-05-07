package config

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
)

func ReadConfig() (Configuration, error) {
	pwd, err := os.Getwd()
	fmt.Printf("Current Path is %s\n", pwd)
	file, err := ioutil.ReadFile(pwd + "/src/fault-tolerance/config/config.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		return Configuration{}, err
	}
	fmt.Printf("Content Length %d\n", len(string(file)))
	var configuration Configuration
	err = json.Unmarshal(file, &configuration)
	if err != nil {
		return Configuration{}, err
	}
	return configuration, nil
}
