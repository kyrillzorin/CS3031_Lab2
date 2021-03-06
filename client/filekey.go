package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// File Key Struct
type FileKey struct {
	Id    string
	User  string
	Owner string
	Name  string
	Key   []byte
}

// Create New File Key
func NewFileKey(user string, owner string, name string, key []byte) *FileKey {
	f := new(FileKey)
	f.User = user
	f.Owner = owner
	f.Name = name
	f.Key = key
	return f
}

// Share a file key on server
func (f *FileKey) Share() error {
	message, err := json.Marshal(f)
	if err != nil {
		return err
	}
	// Sign the request
	signature, err := sign(ClientPrivateKey, message)
	if err != nil {
		return err
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(SignedRequest{message, signature})
	res, err := http.Post(Server+"/sharefile", "application/json; charset=utf-8", b)
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

// Revoke a file key on server
func (f *FileKey) Revoke() error {
	message, err := json.Marshal(f)
	if err != nil {
		return err
	}
	// Sign the request
	signature, err := sign(ClientPrivateKey, message)
	if err != nil {
		return err
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(SignedRequest{message, signature})
	res, err := http.Post(Server+"/revokefile", "application/json; charset=utf-8", b)
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

// Get a file key from server
func GetFileKey(owner string, filename string) (filekey *FileKey, err error) {
	res, err := http.Get(Server + "/users/" + owner + "/" + filename + "/key/" + ClientUser)
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
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&filekey)
	return
}
