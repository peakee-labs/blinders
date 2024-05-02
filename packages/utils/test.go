package utils

import (
	"os"
	"testing"
)

// SkipTestOnEnvironment check if current environments match env.
// If true, test will be skip
func SkipTestOnEvironment(t *testing.T, env string) {
	if env == os.Getenv("ENVIRONMENT") {
		t.SkipNow()
	}
}
