package utils

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"
)

type StitcherTokenPayload struct {
	Sub             string   `json:"sub"`
	Aud             string   `json:"aud"`
	CognitoGroups   []string `json:"cognito:groups"`
	EmailVerified   bool     `json:"email_verified"`
	AuthTime        int      `json:"auth_time"`
	Iss             string   `json:"iss"`
	CognitoUsername string   `json:"cognito:username"`
	Exp             int      `json:"exp"`
	Iat             int      `json:"iat"`
	Email			string	 `json:"email"`
}

func TestReadCredentials(t *testing.T) {
	testCreds, err := ReadCredentials("credentials_test.json")
	if err != nil {
		t.Errorf("Error reading creds file")
	}
	if testCreds.SessionToken != "abcd1234" {
		t.Errorf("Session token wrong. Expected: %s , got: %s", "abcd1234", testCreds.SessionToken)
	}

	_, err = ReadCredentials("credentials_test.txt")
	if err == nil {
		t.Error("ReadCredentials didn't return an error when reading invalid file")
	}

	_, err = ReadCredentials("askdfjasdkjfaksjfsjakdf")
	if err == nil {
		t.Error("ReadCredentials didn't return an error when reading non-existent file")
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

	jwtMissingEmail := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhc2RmYXNkZiIsImF1ZCI6ImFzZGZhc2RmIiwiY29nbml0bzpncm91cHMiOlsiUHJlbWl1bSJdLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXV0aF90aW1lIjoxNjA5ODE0NTcwLCJpc3MiOiJodHRwczovL2NvZ25pdG8taWRwLnVzLWVhc3QtMS5hbWF6b25hd3MuY29tL3VzLWVhc3QtMV9hc2RmIiwiY29nbml0bzp1c2VybmFtZSI6ImFzZGYxMjM0IiwiZXhwIjoxNjA5ODkyMTQ0LCJpYXQiOjE2MDk4ODg1NDV9.e1ssVLtMv8ZO2KZ2KHdCTVQlNNUIbNXGDUzI2a2s138"
	missingEmailIsValid, missingEmailIsValidReason := IsStitcherTokenValid(jwtMissingEmail)
	if missingEmailIsValid || missingEmailIsValidReason != "JWT missing fields" {
		t.Errorf("JWT with email missing was validated or reason was wrong: '%s'", missingEmailIsValidReason)
	}

	j, err := json.Marshal(StitcherTokenPayload{
		Aud: "",
		Sub: "",
		CognitoGroups: []string{"Premium"},
		EmailVerified: true,
		AuthTime: 1609814570,
		Iss: "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_asdf",
		CognitoUsername: "asdf1234",
		Exp: 1000000000, // Way in the past
		Iat: 1609888545,
		Email: "asdf@asdf.lol",
	})
	if err != nil {
		t.Error(err)
	}

	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." + base64.URLEncoding.EncodeToString(j) +
		".Z41mfxmgR0CsdXs_UTqddKMSwNlBqINHmxOWfdT7Vng"

	expiredTokenIsValid, expiredTokenIsValidReason := IsStitcherTokenValid(expiredToken)
	if expiredTokenIsValid || expiredTokenIsValidReason != "token is expired"  {
		t.Errorf("Expired JWT was validated or reason was wrong: '%s'", expiredTokenIsValidReason)
	}

	j2, err := json.Marshal(StitcherTokenPayload{
		Aud: "",
		Sub: "",
		CognitoGroups: []string{"Premium"},
		EmailVerified: true,
		AuthTime: 1609814570,
		Iss: "https://cognito-idp.us-east-1.amazonaws.com/us-east-1_asdf",
		CognitoUsername: "asdf1234",
		Exp: int(time.Now().Unix() + 300),   // 5 minutes from now
		Iat: 1609888545,
		Email: "asdf@asdf.lol",
	})
	if err != nil {
		t.Error(err)
	}

	goodToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." + base64.URLEncoding.EncodeToString(j2) +
				".Z41mfxmgR0CsdXs_UTqddKMSwNlBqINHmxOWfdT7Vng"

	goodJWTisValid, goodJWTisValidReason := IsStitcherTokenValid(goodToken)
	if !goodJWTisValid  {
		t.Errorf("Valid token was not validated: '%s'", goodJWTisValidReason)
	}
}