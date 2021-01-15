package utils

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"
)

type Credentials struct {
	SessionToken string	`json:"session_token"`
	StitcherNewToken string `json:"stitcher_new_token"`
}

func ReadCredentials(file string) (Credentials, error) {
	var creds Credentials
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return creds, err
	}
	err = json.Unmarshal(data, &creds)
	if err != nil {
		return creds, err
	}
	return creds, nil
}

func IsStitcherTokenValid(jwt string) (bool, string) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return false, "invalid JWT format"
	}

	// Add padding if needed
	segment := parts[1]
	if l := len(segment) % 4; l > 0 {
		segment += strings.Repeat("=", 4-l)
	}

	payloadString, e := base64.URLEncoding.DecodeString(segment)
	if e != nil {
		return false, "invalid JWT format"
	}

	var payload map[string]interface{}
	err := json.Unmarshal(payloadString, &payload)
	if err != nil {
		return false, "unable to unmarshal to JSON"
	}

	_, emailFound := payload["email"]
	_, cognitoUsernameFound := payload["cognito:username"]
	exp, expFound := payload["exp"]
	if !emailFound || !cognitoUsernameFound || !expFound {
		return false, "JWT missing fields"
	}

	expirationFloat := exp.(float64)
	now := time.Now().Unix()

	if int64(expirationFloat) <= now {
		return false, "token is expired"
	}

	return true, ""
}