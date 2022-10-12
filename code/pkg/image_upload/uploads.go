package image_upload

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"maryam/image-transcode/pkg/config_reader"
	"maryam/image-transcode/pkg/image_scaling"
	"math"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func VerifyImageConstraints(image *[]byte, config *config_reader.Config) (*image_scaling.ScalingImage, error) {

	log.Printf("Body size is %v bytes\n", len(*image))
	if len(*image) == 0 {
		return nil, errors.New("body is empty")
	}
	if len(*image) > int(config.S3ObjectMaxSizeInMb)*int(math.Pow10(6)) {
		return nil, errors.New("body size is too large")
	}
	sI, err := image_scaling.NewImage(bytes.NewReader(*image))
	if err != nil {
		return nil, fmt.Errorf("%v: error reading image data", err)
	}
	var imageTypeOk bool = false
	for _, imageType := range config.ImageExtensions {
		if imageType == sI.GetFormat() {
			imageTypeOk = true
		}
	}
	if !imageTypeOk {
		return nil, errors.New("image type is not allowed")
	}

	return sI, nil
}

func UploadImageObject(ctx *context.Context, client *s3.Client, bucketName string, objectKey string, imageInfo *image_scaling.ScalingImage, object *[]byte) (err error) {
	log.Printf("Uploading new file")
	uploader := manager.NewUploader(client)
	_, err = uploader.Upload(*ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		Body:        bytes.NewReader(*object),
		ContentType: aws.String(fmt.Sprintf("image/%s", imageInfo.GetFormat())),
	})
	if err != nil {
		return fmt.Errorf("%v: failed to upload image to destination bucket %s", err, bucketName)
	}
	log.Printf("Completed uploading to s3")
	return nil
}
