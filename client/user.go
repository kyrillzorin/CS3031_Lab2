package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
)

type User struct {
	Id       string
	Username string
	PubKey   *rsa.PublicKey
}

func NewUser(username string, pubkey *rsa.PublicKey) *User {
	u := new(User)
	u.Username = username
	u.PubKey = pubkey
	return u
}

func (u *User) Register() error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	res, err := http.Post(Server+"/register", "application/json; charset=utf-8", b)
	if err != nil {
		return err
	}
	if res.Body == nil {
		return errors.New("Empty Response")
	}
	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}
	if response.Status != "success" {
		return errors.New(response.Error)
	}
	return err
}

func GetUser(username string) (user *User, err error) {
	res, err := http.Get(Server + "/users/" + username)
	if err != nil {
		return
	}
	if res.Body == nil {
		err = errors.New("Empty Response")
		return
	}
	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return
	}
	if response.Status == "failure" {
		err = errors.New(response.Error)
		return
	}
	err = json.NewDecoder(res.Body).Decode(&user)
	return
}
