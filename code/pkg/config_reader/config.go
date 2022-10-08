// Package config_reader provides a way to read configurations specific to the project into Config structure
package config_reader

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/magiconair/properties"
)

//Config represents the data structure read from proprties file
//some defaults are provided, while others like paths should be set explicitely
type Config struct {
	DestResolutionY         int    `properties:"resolutionY,default=200"`
	DestResolutionX         int    `properties:"resolutionX,default=200"`
	SourceFilePath          string `properties:"sourceFilePath"`
	DestinationFilePath     string `properties:"destinationFilePath"`
	SourceS3BucketName      string `properties:"sourceS3Bucket"`
	DestinationS3BucketName string `properties:"destinationS3Bucket"`
<<<<<<< Updated upstream
=======
	S3ObjectMaxSizeInMb     int64  `properties:"s3ObjectMaxSizeInMb,default=100"`
>>>>>>> Stashed changes
}

//GetConfigsFromDir reads application properties from a given root directory
func GetConfigsFromDir(configRootDir string) (config *Config, err error) {
	defer func() {
		if err_panic := recover(); err_panic != nil {
			err = fmt.Errorf("%v: failed to get config file", err_panic)
		}
	}()
	log.Printf("Reading configurations from properties files")

	var propertiesFiles = []string{path.Join(configRootDir, "application.properties")}
	env, ok := os.LookupEnv("ENV")
	if ok {
		log.Printf("Adding property file for %s environment", env)
		propertiesFiles = append(propertiesFiles, path.Join(configRootDir, fmt.Sprintf("application-%s.properties", env)))
	}

	p, err := properties.LoadFiles(propertiesFiles, properties.UTF8, true)
	if err != nil {
		return nil, fmt.Errorf("%v: unable to load config file", err)
	}
	config = &Config{}
	err = p.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("%v: unable to decode config file to struct", err)
	}
	log.Printf("Loaded configurations")
	return config, nil
}

//GetConfigs is a wrapper around GetConfigsFromDir
//By default it assumes that configs are in local ./configs/ directory
func GetConfigs() (config *Config, err error) {
	defer func() {
		if err_panic := recover(); err_panic != nil {
			err = fmt.Errorf("%v: failed to get config file", err_panic)
		}
	}()
	return GetConfigsFromDir("configs/")
}
