package utils

import (
	"encoding/json"
	"fmt"
)

type Cred struct {
	Type     string
	User     string
	Password string
}

func Simple() (string, error) {
	var user, password string
	fmt.Print("User: ")
	_, err := fmt.Scanln(&user)
	if err != nil {
		return "", err
	}
	fmt.Print("Password: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		return "", err
	}
	return MakeCred("simple", user, password)
}

func MakeCred(authType string, user string, password string) (string, error) {
	cred := Cred{
		Type:     authType,
		User:     user,
		Password: password,
	}
	buf, err := json.Marshal(cred)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func GetCred(auth string) (*Cred, error) {
	var buf Cred
	err := json.Unmarshal([]byte(auth), &buf)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
