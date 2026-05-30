package main

import (
	"log"
	"os"

	"github.com/Dedushka-Lenin/DockerControllerGolang/internal/adapters/config"
)

const configPath = "config/config.json"

func launchApp() error {
	_, err := os.Stat(configPath)
	if err != nil {
		return err
	}

	conf, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	server, err := NewServer(conf)
	if err != nil {
		return err
	}

	err = server.Run()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := launchApp()

	if err != nil {
		log.Fatal(err)
	}
}
