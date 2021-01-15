package utils

import (
	"testing"
)

func TestIsValidURL(t *testing.T) {
	if !IsValidURL("https://my.site/podcast") {
		t.Errorf("Valid URL did not validate")
	}

	if IsValidURL("asdf;lkj") {
		t.Errorf("Invalid URL validated")
	}
}