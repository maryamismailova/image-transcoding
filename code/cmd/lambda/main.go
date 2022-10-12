package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"maryam/image-transcode/pkg/config_reader"
	"maryam/image-transcode/pkg/image_scaling"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func downloadObject(ctx *context.Context, client *s3.Client, bucketName string, objectKey string, writer io.WriterAt) (err error) {

	downloader := manager.NewDownloader(client)
	log.Printf("Download object %s from bucket %s \n", objectKey, bucketName)
	_, err = downloader.Download(*ctx, writer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("%v: error download object from s3bucket to temporary file", err)
	}
	return nil
}

func uploadObject(ctx *context.Context, client *s3.Client, bucketName string, objectKey string, reader io.Reader, imageInfo *image_scaling.ScalingImage) (err error) {

	log.Printf("Uploading new file")
	uploader := manager.NewUploader(client)
	_, err = uploader.Upload(*ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        reader,
		ContentType: aws.String(fmt.Sprintf("image/%s", imageInfo.GetFormat())),
	})
	if err != nil {
		return fmt.Errorf("%v: failed to upload scaled image to destination bucket %s", err, bucketName)
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

			// download file into memory
			var buff []byte = make([]byte, 1024)
			writer := manager.NewWriteAtBuffer(buff)
			err = downloadObject(&ctx, client, bucketName, objectKey, writer)
			if err != nil {
				return fmt.Errorf("%v: error downloading object", err)
			}

			// Start scaling process
			image_old_reader := bufio.NewReader(bytes.NewReader(writer.Bytes()))
			image_new_readwriter := bufio.NewReadWriter(bufio.NewReader(bytes.NewBuffer(make([]byte, 1024))), bufio.NewWriter(bytes.NewBuffer(make([]byte, 1024))))
			image, err := image_scaling.ScaleImage(image_old_reader, image_new_readwriter, transcoding.GetResolutionY(), transcoding.GetResolutionX())
			if err != nil {
				return fmt.Errorf("%v: error scaling image", err)
			}

			// Upload file
			err = uploadObject(&ctx, client, conf.DestinationS3BucketName, destObjectKey, image_new_readwriter, image)
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
