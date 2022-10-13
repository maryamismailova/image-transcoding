# Cloud Image Transcoding

## Description:

Get image transcoded to fit in the square with resolutions of your choice once it gets uploaded to S3 bucket.

## How to use:

### Within SAM application template with automated pipeline

Check out the list of parameters in [template.yaml](code/template.yaml) and for those that you need to update, modify the override-parameters in [samconfig.toml](code/samconfig.toml)

Currently samconfig.toml is configured for 3 different enviornments. Modify parameters accordingly.

### Parameters described

```
  // Application specific parameters
  ImageMaxAllowedSize: Maximum allowed size for input images. Others won't be transcoded.

  ImageMaxAllowedResolutions: List of square resolutions to which input images will be scaled. Should be provided in form of: "Y1xX2;Y2xX2;..."

  // Resource related parameters
  SourceBucketName: bucket name for s3 bucket to upload source images to. Must be created by sam only. Otherwise deployment will fail

  DestinationBucketName: bucket name for the upload of transcoded images. Same as Source bucket, should only be created by sam deploy.

  Env: Environment prefix

  APIFunctionNamePrefix: Name for api lambda (DRAFT)

  MemorySize: Allowed memory size for lambda (currently shared across all lambda)

  FunctionNamePrefix: Name prefix of image resizer lambda function. Will be added ${Env} value as suffix.

  // Additional parameters to improve linking between function and the source code from which it was created
  // This variables are set by pipeline itself/
  WorkflowRun:

  CommitId:

  Branch:

```

## Configuring the function itself (running from binary)

One way is to pass configuration properties through a file, another one is by environment variables. Properties file is preferred and more stable one. Environment variables were introduced due to ease of integration with lambda and SAM stacks.. SAM sucks.

### Configuration properties

Pass configs to code/configs/ directory by default. Though it can be overriden.

Allowed content of application.properties is defined by [Config type](code/pkg/config_reader/config.go) in configurations package.

Example:

```
sourceFilePath=/tmp
destinationFilePath=/tmp
sourceS3Bucket=s3-input
destinationS3Bucket=s3-dest
destinationFilePath=outputs/test-1.png
sourceS3Bucket=source
destinationS3Bucket=destination
transcodingResolutions = 1024x1024;2048x2048
imageExtensions=png;jpeg
```

Can be overriden with application-<ENV>.properties files if ENV enviornment variable is defined.

### Environment variables

Some configurations can be overriden by means of environment variables. Thos are:

```
S3_SOURCE_BUCKET  ->  which maps to sourceS3Bucket

S3_DESTINATION_BUCKET ->  which maps to destinationS3Bucket

IMAGE_MAX_ALLOWED_SIZE -> which maps to s3ObjectMaxSizeInMb

IMAGE_ALLOWED_RESOLUTIONS -> which maps to transcodingResolutions

IMAGE_ALLOWED_EXTENSIONS -> which maps to imageExtensions
```

## Configure pipeline

For pipeline to run please create a user in AWS, give the user permissions to resources defined within the stack(Shall be updated..). And defined AWS access key and key id and region in repo's secrets.
