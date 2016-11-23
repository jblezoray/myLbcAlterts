package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	SearchTerms string
	SMTPUser    SMTPUser
}

func ReadConfigFile(filename string) (Configuration, error) {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	return configuration, err
}
