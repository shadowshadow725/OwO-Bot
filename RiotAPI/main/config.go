package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Keys struct {
		Discord string `yaml:"discordKey"`
		Riot string `yaml:"riotKey"`
	} `yaml:"Keys"`
	Users struct {
		DiscordID string `yaml:"discordID"`
		IGN string `yaml:"IGN"`
	} `yaml:"Users"`
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
func ParseConfig(filename string, cfg *Config){

	f, err := os.Open(filename)
	if err != nil {
		processError(err)

	}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)

	if err != nil {
		processError(err)
	}

	//fmt.Printf("%s\n ",cfg.Users.IGN)
	fmt.Printf("%+v", cfg)
	f.Close()

	return
}
