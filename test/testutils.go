package test

import (
	"reflect"
	"testing"
)

func AssertString(t *testing.T, valueName string, expected string, got string) {
	if expected != got {
		t.Errorf("%s wrong - expected: '%s', got: '%s'", valueName, expected, got)
	}
}

func AssertTypesAreEqual(t *testing.T, first interface{}, second interface{}) {
	firstType := reflect.TypeOf(first)
	secondType := reflect.TypeOf(second)

	if firstType != secondType {
		t.Errorf("%s and %s are not the same type", firstType, secondType)
	}
}

