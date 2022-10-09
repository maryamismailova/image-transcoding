package main

import (
	"log"
	"maryam/image-transcode/pkg/config_reader"
	"maryam/image-transcode/pkg/image_scaling"
)

func TestImageScaling() {
	config, err := config_reader.GetConfigs()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	log.Printf("Loaded config: %+v\n", config)
	log.Printf("%v\n", len(config.TranscodingResolutions))
	for _, transcoding := range config.TranscodingResolutions {
		_, err = image_scaling.ScaleImageFromSource(config.SourceFilePath, config.DestinationFilePath, transcoding.GetResolutionY(), transcoding.GetResolutionX())
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
	}
}

func main() {
	TestImageScaling()
}
