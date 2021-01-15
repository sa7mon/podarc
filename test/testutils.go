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

func AssertEqual(t *testing.T, first interface{}, second interface{}) {
	if first != second {
		t.Errorf("'%s' and '%s' are not equal", first, second)
	}
}

func AssertNotEmpty(t *testing.T, testName string, s string) {
	if len(s) < 1 {
		t.Errorf("[%s] string is empty", testName)
	}
}