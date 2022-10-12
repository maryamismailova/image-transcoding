package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"maryam/image-transcode/pkg/config_reader"
	"maryam/image-transcode/pkg/image_upload"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func HandleImageUploadRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	log.Println("Image upload request received")
	log.Printf("Request: rawpath=%v, rawquery=%v, httpobject=%+v\n", request.RawPath, request.RawQueryString, request.RequestContext.HTTP)
	var err error = nil
	if request.RequestContext.HTTP.Method != "POST" {
		return events.LambdaFunctionURLResponse{
			StatusCode:      http.StatusMethodNotAllowed,
			Body:            "Method not alowed",
			IsBase64Encoded: false,
		}, nil
	}
	if request.RequestContext.HTTP.Path != "/upload" {
		return events.LambdaFunctionURLResponse{
			StatusCode:      http.StatusNotAcceptable,
			Body:            "No such request path defined",
			IsBase64Encoded: false,
		}, nil
	}
	//Get image name if provided
	var image_name_prefix = request.RequestContext.RequestID
	if len(request.RawQueryString) != 0 {
		re := regexp.MustCompile(`name=(.+)`)
		n := re.FindStringSubmatch(request.RawQueryString)
		if len(n) != 0 {
			image_name_prefix = n[1]
		}
	}

	log.Println("Loding configurations")
	configs, err := config_reader.GetConfigs()
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: 500}, err
	}
	log.Println("Loded configurations")

	log.Println("Decoding image data from base64")
	image_bytes, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest}, nil
	}
	log.Println("Decoded image")

	// Verify image
	log.Println("Verifying whether image corresponds to requirements")
	image_obj, err := image_upload.VerifyImageConstraints(&image_bytes, configs)
	if err != nil {
		log.Printf("Bad request event: %+v\n", events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest})
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusBadRequest}, nil
	}
	log.Println("Image verification complete")

	log.Printf("Image received of format %s\n", image_obj.GetFormat())

	// Image upload
	log.Println("Started uploading to S3")
	image_name_full := fmt.Sprintf("%s.%s", image_name_prefix, image_obj.GetFormat())
	log.Println("Initializing aws session")
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError}, err
	}
	client := s3.NewFromConfig(cfg)
	log.Println("Session initialized")
	log.Printf("Uploading object %s to bucket %s", image_name_full, configs.SourceS3BucketName)
	err = image_upload.UploadImageObject(&ctx, client, configs.SourceS3BucketName, image_name_full, image_obj, &image_bytes)
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusInternalServerError}, err
	}
	log.Println("Finished uploading to S3")

	return events.LambdaFunctionURLResponse{
		StatusCode: http.StatusOK,
		Body:       fmt.Sprintf("Image was uploaded to %s/%s", configs.SourceS3BucketName, image_name_full),
	}, nil

}

func main() {
	lambda.Start(HandleImageUploadRequest)
}
