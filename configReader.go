package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	SearchTerms      string
	DatabaseFilepath string
	SMTPUser         SMTPUser
	MailFrom         string
	MailTo           string
}

func ReadConfigFile(filename string) (Configuration, error) {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	return configuration, err
}
