package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type Credentials struct {
	SessionToken string	`json:"session_token"`
	StitcherNewToken string `json:"stitcher_new_token"`
}

func ReadCredentials(file string) Credentials {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Print(err)
	}
	var creds Credentials
	err = json.Unmarshal(data, &creds)
	if err != nil {
		fmt.Println("Error reading creds file: ", err)
	}
	return creds
}

func IsStitcherTokenValid(jwt string) (bool, string) {
	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return false, "invalid JWT format"
	}

	payload_string, e := base64.StdEncoding.DecodeString(parts[1])
	if e != nil {
		return false, "invalid JWT format"
	}

	var payload map[string]interface{}
	err := json.Unmarshal(payload_string, &payload)
	if err != nil {
		return false, "unable to unmarshal to JSON"
	}

	_, email_found := payload["email"]
	_, cognito_username_found := payload["cognito:username"]
	exp, exp_found := payload["exp"]
	if !email_found || !cognito_username_found || !exp_found {
		return false, "JWT missing fields"
	}

	expiration_float := exp.(float64)
	now := time.Now().Unix()

	if int64(expiration_float) <= now {
		return false, "token is expired"
	}

	return true, ""
}