package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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