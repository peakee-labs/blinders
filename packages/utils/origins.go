package utils

import (
	"fmt"
	"os"
)

func GetOriginsFromEnv() string {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "prod" {
		return "https://app.peakee.co"
	} else if environment == "staging" {
		return fmt.Sprintf("https://%s.app.peakee.co", environment)
	} else if environment == "dev" {
		return "*"
	}

	panic("unknown deployment stage")
}
