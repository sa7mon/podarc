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
	tooShortisValid, tooShortisValidReason := IsStitcherTokenValid(tooShort)
	if tooShortisValid || tooShortisValidReason != "invalid JWT format" {
		t.Errorf("too short token was validated or reason was wrong: '%s'", tooShortisValidReason)
	}

	badBase64 := "!!asdf.!!!asdf.asdf!!!"
	badBase64isValid, badBase64isValidReason := IsStitcherTokenValid(badBase64)
	if badBase64isValid || badBase64isValidReason != "invalid JWT format" {
		t.Errorf("bad Base64 token was validated or reason was wrong: '%s'", badBase64isValidReason)
	}

	invalidJSONBase64 := "YXNkZmFzZGY=.YXNkZmFzZGY=.YXNkZmFzZGY="
	invalidJSONBase64isValid, invalidJSONBase64isValidReason := IsStitcherTokenValid(invalidJSONBase64)
	if invalidJSONBase64isValid || invalidJSONBase64isValidReason != "unable to unmarshal to JSON" {
		t.Errorf("invalid JSON Base64 token was validated or reason was wrong: '%s'", invalidJSONBase64isValidReason)
	}

	jwtMissingEmail := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhc2RmYXNkZiIsImF1ZCI6ImFzZGZhc2RmIiwiY29nbml0bzpncm91cHMiOlsiUHJlbWl1bSJdLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXV0aF90aW1lIjoxNjA5ODE0NTcwLCJpc3MiOiJodHRwczovL2NvZ25pdG8taWRwLnVzLWVhc3QtMS5hbWF6b25hd3MuY29tL3VzLWVhc3QtMV9hc2RmIiwiY29nbml0bzp1c2VybmFtZSI6ImFzZGYxMjM0IiwiZXhwIjoxNjA5ODkyMTQ0LCJpYXQiOjE2MDk4ODg1NDUsImVtYWlsIjoiYXNkZkBlbWFpbC5jb20ifQ.Z41mfxmgR0CsdXs_UTqddKMSwNlBqINHmxOWfdT7Vng"
	missingEmailIsValid, missingEmailIsValidReason := IsStitcherTokenValid(jwtMissingEmail)
	if missingEmailIsValid || missingEmailIsValidReason != "JWT missing fields" {
		t.Errorf("JWT with email missing was validated or reason was wrong: '%s'", missingEmailIsValidReason)
	}

	goodJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhc2RmYXNkZiIsImF1ZCI6ImFzZGZhc2RmIiwiY29nbml0bzpncm91cHMiOlsiUHJlbWl1bSJdLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXV0aF90aW1lIjoxNjA5ODE0NTcwLCJpc3MiOiJodHRwczovL2NvZ25pdG8taWRwLnVzLWVhc3QtMS5hbWF6b25hd3MuY29tL3VzLWVhc3QtMV9hc2RmIiwiY29nbml0bzp1c2VybmFtZSI6ImFzZGYxMjM0IiwiZXhwIjoxNjA5ODkyMTQ0LCJpYXQiOjE2MDk4ODg1NDUsImVtYWlsIjoiYXNkZkBlbWFpbC5jb20ifQ.Z41mfxmgR0CsdXs_UTqddKMSwNlBqINHmxOWfdT7Vng"
	goodJWTisValid, goodJWTisValidReason := IsStitcherTokenValid(goodJWT)
	if !goodJWTisValid  {
		t.Errorf("Valid token was not validated: '%s'", goodJWTisValidReason)
	}
}