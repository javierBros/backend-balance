package application

import "os"

type EnvironmentVariables struct {
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	SESRegion          string
	S3Bucket           string
	DestinationEmail   string
	SourceEmail        string
}

func NewEnvironmentVariables() *EnvironmentVariables {
	envVariables := &EnvironmentVariables{
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SESRegion:          os.Getenv("SES_REGION"),
		S3Bucket:           os.Getenv("S3_BUCKET"),
		DestinationEmail:   os.Getenv("DESTINATION_EMAIL"),
		SourceEmail:        os.Getenv("SOURCE_EMAIL"),
	}

	if envVariables.AWSAccessKeyID == "" {
		envVariables.AWSAccessKeyID = "defaultAccessKey"
	}
	if envVariables.AWSSecretAccessKey == "" {
		envVariables.AWSSecretAccessKey = "defaultSecretKey"
	}
	if envVariables.SESRegion == "" {
		envVariables.SESRegion = "us-east-1"
	}
	if envVariables.S3Bucket == "" {
		envVariables.S3Bucket = "defaultS3Bucket"
	}
	if envVariables.DestinationEmail == "" {
		envVariables.DestinationEmail = "javiteck031@gmail.com"
	}
	if envVariables.SourceEmail == "" {
		envVariables.SourceEmail = "javiteck031@gmail.com"
	}

	return envVariables
}
