package main

import (
	"context"
	"fmt"
	"maryam/image-transcode/pkg/config_reader"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleS3Trigger(ctx context.Context, s3Event events.S3Event) (err error) {
	config, err := config_reader.GetConfigs()
	if err != nil {
		return fmt.Errorf("%v: unable to read configs", err)
	}
	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)
		if s3.Object.Size/1024 >= config.S3ObjectMaxSizeInMb {
			return fmt.Errorf("%v: max size for s3 object exceeded", err)
		}
	}
	return err
}

func main() {
	lambda.Start(HandleS3Trigger)

}
