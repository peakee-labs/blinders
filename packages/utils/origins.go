package utils

import (
	"fmt"
	"os"
)

func GetOriginsFromEnv() string {
	deploymentStage := os.Getenv("DEPLOYMENT_STAGE")
	if deploymentStage == "prod" {
		return "https://app.peakee.co"
	} else if deploymentStage == "staging" {
		return fmt.Sprintf("https://%s.app.peakee.co", deploymentStage)
	} else if deploymentStage == "dev" {
		return "*"
	}

	panic("unknown deployment stage")
}
