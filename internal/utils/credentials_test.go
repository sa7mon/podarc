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

func TestIsStitcherTokenValid(t *testing.T) {
	tooShort := "asdfasdf1234"
	isValid, isValidReason := IsStitcherTokenValid(tooShort)
	if isValid || isValidReason != "invalid JWT format" {
		t.Errorf("too short token was validated or reason was wrong: '%s'", isValidReason)
	}

	badBase64 := "!!asdf.!!!asdf.asdf!!!"
	isValid, isValidReason = IsStitcherTokenValid(badBase64)
	if isValid || isValidReason != "invalid JWT format" {
		t.Errorf("bad Base64 token was validated or reason was wrong: '%s'", isValidReason)
	}

	invalidJSONBase64 := "YXNkZmFzZGY=.YXNkZmFzZGY=.YXNkZmFzZGY="
	isValid, isValidReason = IsStitcherTokenValid(invalidJSONBase64)
	if isValid || isValidReason != "unable to unmarshal to JSON" {
		t.Errorf("invalid JSON Base64 token was validated or reason was wrong: '%s'", isValidReason)
	}

	jwtMissingEmail := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhc2RmYXNkZiIsImF1ZCI6ImFzZGZhc2RmIiwiY29nbml0bzpncm91cHMiOlsiUHJlbWl1bSJdLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXV0aF90aW1lIjoxNjA5ODE0NTcwLCJpc3MiOiJodHRwczovL2NvZ25pdG8taWRwLnVzLWVhc3QtMS5hbWF6b25hd3MuY29tL3VzLWVhc3QtMV9hc2RmIiwiY29nbml0bzp1c2VybmFtZSI6ImFzZGYxMjM0IiwiZXhwIjoxNjA5ODkyMTQ0LCJpYXQiOjE2MDk4ODg1NDUsImVtYWlsIjoiYXNkZkBlbWFpbC5jb20ifQ.Z41mfxmgR0CsdXs_UTqddKMSwNlBqINHmxOWfdT7Vng"
	isValid, isValidReason = IsStitcherTokenValid(jwtMissingEmail)
	if isValid || isValidReason != "JWT missing fields" {
		t.Errorf("JWT with email missing was validated or reason was wrong: '%s'", isValidReason)
	}

	goodJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhc2RmYXNkZiIsImF1ZCI6ImFzZGZhc2RmIiwiY29nbml0bzpncm91cHMiOlsiUHJlbWl1bSJdLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXV0aF90aW1lIjoxNjA5ODE0NTcwLCJpc3MiOiJodHRwczovL2NvZ25pdG8taWRwLnVzLWVhc3QtMS5hbWF6b25hd3MuY29tL3VzLWVhc3QtMV9hc2RmIiwiY29nbml0bzp1c2VybmFtZSI6ImFzZGYxMjM0IiwiZXhwIjoxNjA5ODkyMTQ0LCJpYXQiOjE2MDk4ODg1NDUsImVtYWlsIjoiYXNkZkBlbWFpbC5jb20ifQ.Z41mfxmgR0CsdXs_UTqddKMSwNlBqINHmxOWfdT7Vng"
	isValid, isValidReason = IsStitcherTokenValid(goodJWT)
	if !isValid  {
		t.Errorf("Valid token was not validated: '%s'", isValidReason)
	}
}