package main

import (
	"fmt"
	"log"
	"maryam/image-transcode/pkg/config_reader"
	"maryam/image-transcode/pkg/image_scaling"
	"os"
	"path/filepath"
)

func TestImageScaling() {
	config, err := config_reader.GetConfigs()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	log.Printf("Loaded config: %+v\n", config)

	for _, transcoding := range config.TranscodingResolutions {
		dir, file := filepath.Split(config.DestinationFilePath)
		file_path := filepath.Join(dir, fmt.Sprintf("%dx%d-%s", transcoding.GetResolutionY(), transcoding.GetResolutionX(), file))
		_, err = image_scaling.ScaleImageFromSource(config.SourceFilePath, file_path, transcoding.GetResolutionY(), transcoding.GetResolutionX())
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
	}
}

func TestImageScalingWithIo() {
	config, err := config_reader.GetConfigs()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	log.Printf("Loaded config: %+v\n", config)

	for _, transcoding := range config.TranscodingResolutions {
		log.Printf("Starting scaling image for %dx%d\n", transcoding.GetResolutionY(), transcoding.GetResolutionX())
		fSrc, err := os.Open(config.SourceFilePath)
		if err != nil {
			log.Fatalf("%v: failed to open source file %s", err, config.SourceFilePath)
		}
		defer fSrc.Close()
		dir, file := filepath.Split(config.DestinationFilePath)
		file_path := filepath.Join(dir, fmt.Sprintf("%dx%d-%s", transcoding.GetResolutionY(), transcoding.GetResolutionX(), file))
		fDest, err := os.Create(file_path)
		if err != nil {
			log.Fatalf("%v: failed to create destination file %s", err, file_path)
		}
		defer fDest.Close()
		_, err = image_scaling.ScaleImage(fSrc, fDest, transcoding.GetResolutionY(), transcoding.GetResolutionX())
		if err != nil {
			log.Fatalf("%v: failed to scale image", err)
		}
		log.Printf("Finished scaling image for %dx%d\n", transcoding.GetResolutionY(), transcoding.GetResolutionX())

	}
}

func main() {
	// TestImageScalingWithIo()
	TestImageScaling()
}
