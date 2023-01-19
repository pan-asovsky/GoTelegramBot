package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"syscall"
)

type Config struct {
	Telegram struct {
		BotToken    string `yaml:"bot_token"`
		WebhookPath string `yaml:"webhook_path"`
		ServerUrl   string `yaml:"server_url"`
	} `yaml:"telegram"`
}

func getConfig() Config {

	filename := "config.yml"
	file := openFile(filename)

	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("Error closing config file: ", err)
		}
	}()

	decoder := yaml.NewDecoder(file)
	config := Config{}

	err := decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Error decoding config file: ", err)
	}

	return config
}

func openFile(name string) *os.File {

	file, err := os.Open(name)
	handleFileOpenErr(err)
	return file
}

func handleFileOpenErr(err error) {

	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalln("Config file does not exist")
		} else if os.IsPermission(err) {
			log.Fatalln("No permission to read config file")
		} else if err.(*os.PathError).Err == syscall.EISDIR {
			log.Fatalln("Config file is a directory")
		} else {
			log.Fatalln("Error opening config file: ", err)
		}
	}
}
