package sample_tests

import (
	"os"
	"testing"
)

func TestEnvVarAvailable(t *testing.T) {
	got := os.Getenv("TASKRUNNER_TEST_ENV")
	if got == "" {
		t.Skip("TASKRUNNER_TEST_ENV not set")
	}
	if got != "expected" {
		t.Fatalf("TASKRUNNER_TEST_ENV = %q, want %q", got, "expected")
	}
}
