package utils

import (
	"github.com/sa7mon/podarc/internal/utils"
	"testing"
)

func TestReadCredentials(t *testing.T) {
	testCreds := utils.ReadCredentials("testCreds.json")
	if testCreds.SessionToken != "abcd1234" {
		t.Errorf("Session token wrong. Expected: %s , got: %s", "abcd1234", testCreds.SessionToken)
	}
}