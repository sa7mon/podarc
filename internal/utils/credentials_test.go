package utils

import (
	"testing"
)

func TestReadCredentials(t *testing.T) {
	testCreds := ReadCredentials("credentials_test.json")
	if testCreds.SessionToken != "abcd1234" {
		t.Errorf("Session token wrong. Expected: %s , got: %s", "abcd1234", testCreds.SessionToken)
	}
}