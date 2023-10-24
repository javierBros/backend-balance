package application

import "os"

type EnvironmentVariables struct {
	SESRegion        string
	DestinationEmail string
	SourceEmail      string
}

func NewEnvironmentVariables() *EnvironmentVariables {
	envVariables := &EnvironmentVariables{
		SESRegion:        os.Getenv("SES_REGION"),
		DestinationEmail: os.Getenv("DESTINATION_EMAIL"),
		SourceEmail:      os.Getenv("SOURCE_EMAIL"),
	}

	if envVariables.SESRegion == "" {
		envVariables.SESRegion = "us-east-1"
	}
	if envVariables.DestinationEmail == "" {
		envVariables.DestinationEmail = "javiteck031@gmail.com"
	}
	if envVariables.SourceEmail == "" {
		envVariables.SourceEmail = "javiteck031@gmail.com"
	}

	return envVariables
}
