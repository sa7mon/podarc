package testutils

import "testing"

func AssertString(t *testing.T, valueName string, expected string, got string) {
	if expected != got {
		t.Errorf("%s wrong - expected: '%s', got: '%s'", valueName, expected, got)
	}
}
