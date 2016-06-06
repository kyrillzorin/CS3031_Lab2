package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// User Struct
type User struct {
	Id       string
	Username string
	PubKey   *rsa.PublicKey
}

// Create new user
func NewUser(username string, pubkey *rsa.PublicKey) *User {
	u := new(User)
	u.Username = username
	u.PubKey = pubkey
	return u
}

// Register user on server
func (u *User) Register() error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
	res, err := http.Post(Server+"/register", "application/json; charset=utf-8", b)
	if res == nil {
		return errors.New("Empty Response")
	}
	if res.Body == nil {
		return errors.New("Empty Response")
	}
	defer res.Body.Close()
	if err != nil {
		return err
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

// Get a user from server
func GetUser(username string) (user *User, err error) {
	res, err := http.Get(Server + "/users/" + username)
	if res == nil {
		err = errors.New("Empty Response")
		return
	}
	if res.Body == nil {
		err = errors.New("Empty Response")
		return
	}
	defer res.Body.Close()
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var response Response
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&response)
	if err != nil {
		return
	}
	if response.Status == "failure" {
		err = errors.New(response.Error)
		return
	}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&user)
	return
}
