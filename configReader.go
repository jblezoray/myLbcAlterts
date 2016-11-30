package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Configuration struct {
	Searches         []Search
	DatabaseFilepath string
	SMTPUser         SMTPUser
	MailFrom         string
	MailTo           string
}

type Search struct {
	Name  string
	Terms string
}

type SMTPUser struct {
	Username string
	Password string
	Server   string
	Port     int
}

func ReadConfigFile(filename string) (Configuration, error) {
	file, _ := os.Open(filename)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		return configuration, err
	}
	return configuration, checkConfiguration(configuration)
}

func checkConfiguration(config Configuration) error {
	for _, search := range config.Searches {
		if search.Name == "" {
			return errors.New("Each search must have a field 'name'.")
		}
		if search.Terms == "" {
			return errors.New("Each search must have a field 'Terms'.")
		}
	}
	return nil
}
