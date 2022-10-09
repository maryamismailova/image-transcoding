package main

import (
	"context"
	"fmt"
	"log"
	"maryam/image-transcode/pkg/config_reader"
	"maryam/image-transcode/pkg/image_scaling"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func downloadObject(ctx *context.Context, client *s3.Client, bucketName string, objectKey string, tmpObjectFilePath string) (err error) {
	tmpObjectFile, err := os.Create(tmpObjectFilePath)
	if err != nil {
		return fmt.Errorf("%v: error creating file to download object to", err)
	}
	defer tmpObjectFile.Close()

	downloader := manager.NewDownloader(client)
	log.Printf("Download object %s from bucket %s into file %s\n", objectKey, bucketName, tmpObjectFilePath)
	_, err = downloader.Download(*ctx, tmpObjectFile, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("%v: error download object from s3bucket to temporary file", err)
	}
	return nil
}

func uploadObject(ctx *context.Context, client *s3.Client, bucketName string, objectKey string, objectFilePath string, imageInfo *image_scaling.ScalingImage) (err error) {
	objectFile, err := os.Open(objectFilePath)
	if err != nil {
		return fmt.Errorf("%v: error opening file to download", err)
	}
	defer objectFile.Close()

	log.Printf("Uploading new file")
	uploader := manager.NewUploader(client)
	_, err = uploader.Upload(*ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        objectFile,
		ContentType: aws.String(fmt.Sprintf("image/%s", imageInfo.GetFormat())),
	})
	if err != nil {
		return fmt.Errorf("%v: failed to upload scaled image %s to destination bucket %s", err, objectFilePath, bucketName)
	}
	log.Printf("Completed uploading to s3")
	return nil
}

func HandleS3Trigger(ctx context.Context, s3Event events.S3Event) (err error) {
	log.Printf("Received event to scale a new image")

	// init
	log.Printf("Initializing aws session")
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("%v: error establishing session with aws", err)
	}

	client := s3.NewFromConfig(cfg)

	conf, err := config_reader.GetConfigs()
	if err != nil {
		return fmt.Errorf("%v: unable to read configs", err)
	}
	log.Printf("Initialization finished")

	// traverse s3 events
	log.Printf("Traversing s3 event")
	for _, record := range s3Event.Records {
		for _, transcoding := range conf.TranscodingResolutions {
			s3obj := record.S3
			bucketName := s3obj.Bucket.Name
			objectKey := s3obj.Object.Key
			destObjectKey := fmt.Sprintf("%dx%d/%s", transcoding.GetResolutionY(), transcoding.GetResolutionX(), objectKey)
			tmpObjectFilePath := fmt.Sprintf("/tmp/%s", objectKey)
			tmpScaledObjectFilePath := fmt.Sprintf("/tmp/%dx%d-%s", transcoding.GetResolutionY(), transcoding.GetResolutionX(), objectKey)

			// download file into disk
			err = downloadObject(&ctx, client, bucketName, objectKey, tmpObjectFilePath)
			if err != nil {
				return fmt.Errorf("%v: error downloading object", err)
			}

			// Start scaling process
			image, err := image_scaling.ScaleImageFromSource(tmpObjectFilePath, tmpScaledObjectFilePath, transcoding.GetResolutionY(), transcoding.GetResolutionX())
			if err != nil {
				return fmt.Errorf("%v: error scaling image", err)
			}

			// Upload file
			err = uploadObject(&ctx, client, conf.DestinationS3BucketName, destObjectKey, tmpScaledObjectFilePath, image)
			if err != nil {
				return fmt.Errorf("%v: error uploading scaled image", err)
			}

		}
	}
	log.Printf("Traversing s3 event finished")
	return err
}

func main() {
	lambda.Start(HandleS3Trigger)

}
