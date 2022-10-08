package main

import (
	_ "image/jpeg"
	_ "image/png"
	"log"
)

func TestImageScaling() {
	config, err := getConfigs()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	log.Printf("Loaded config: %+v\n", config)
	err = scaleImageFromSource(config.SourceFilePath, config.DestinationFilePath, config.DestResolutionY, config.DestResolutionX)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
}

func main() {
	TestImageScaling()
}
