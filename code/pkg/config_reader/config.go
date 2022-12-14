// Package config_reader provides a way to read configurations specific to the project into Config structure
package config_reader

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/magiconair/properties"
)

//Config represents the data structure read from proprties file
//some defaults are provided, while others like paths should be set explicitely
type Config struct {
	DestResolutionY     int    `properties:"resolutionY,default=200"`
	DestResolutionX     int    `properties:"resolutionX,default=200"`
	SourceFilePath      string `properties:"sourceFilePath,default=/tmp"`
	DestinationFilePath string `properties:"destinationFilePath,default=/tmp"`
	//ENV: S3_SOURCE_BUCKET
	SourceS3BucketName string `properties:"sourceS3Bucket,default=input"`
	//ENV: S3_DESTINATION_BUCKET
	DestinationS3BucketName string `properties:"destinationS3Bucket,default=output"`
	//ENV: IMAGE_MAX_ALLOWED_SIZE
	S3ObjectMaxSizeInMb int64 `properties:"s3ObjectMaxSizeInMb,default=100"`
	//ENV: IMAGE_ALLOWED_RESOLUTIONS
	TranscodingResolutions []Resolution `properties:"transcodingResolutions,default=300x300"`
	//ENV: IMAGE_ALLOWED_EXTENSIONS
	ImageExtensions []string `properties:"imageExtensions,default=png;jpeg"`
}

type Resolution string

func (resolution *Resolution) GetResolutionX() (res int) {
	res, _ = strconv.Atoi(strings.Split(string(*resolution), "x")[1])
	return res
}

func (resolution *Resolution) GetResolutionY() (res int) {
	res, _ = strconv.Atoi(strings.Split(string(*resolution), "x")[0])
	return res
}

func (resolution *Resolution) verify() bool {
	if len(strings.Split(string(*resolution), "x")) != 2 {
		return false
	}
	if strings.Compare(string(*resolution), fmt.Sprintf("%dx%d", resolution.GetResolutionY(), resolution.GetResolutionX())) != 0 {
		return false
	}
	return true
}

func overrideFromEnvs(p *properties.Properties) {
	if env, ok := os.LookupEnv("S3_SOURCE_BUCKET"); ok {
		p.Set("sourceS3Bucket", env)
	}
	if env, ok := os.LookupEnv("S3_DESTINATION_BUCKET"); ok {
		p.Set("destinationS3Bucket", env)
	}
	if env, ok := os.LookupEnv("IMAGE_MAX_ALLOWED_SIZE"); ok {
		p.Set("s3ObjectMaxSizeInMb", env)
	}
	if env, ok := os.LookupEnv("IMAGE_ALLOWED_RESOLUTIONS"); ok {
		p.Set("transcodingResolutions", env)
	}
	if env, ok := os.LookupEnv("IMAGE_ALLOWED_EXTENSIONS"); ok {
		p.Set("imageExtensions", env)
	}
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

	overrideFromEnvs(p)

	config = &Config{}
	err = p.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("%v: unable to decode config file to struct", err)
	}

	// Verify complex formats
	for _, transcoding := range config.TranscodingResolutions {
		if !transcoding.verify() {
			return nil, fmt.Errorf("resolution %s with incorrect format. expected <int>x<int> (like 1024x1024)", string(transcoding))
		}
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
