package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type FileKey struct {
	Id    string
	User  string
	Owner string
	Name  string
	Key   []byte
}

func NewFileKey(user string, owner string, name string, key []byte) *FileKey {
	f := new(FileKey)
	f.User = user
	f.Owner = owner
	f.Name = name
	f.Key = key
	return f
}

func (f *FileKey) Share() error {
	message, err := json.Marshal(f)
	signature, _ := sign(ClientPrivateKey, message)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(SignedRequest{message, signature})
	res, err := http.Post(Server+"/sharefile", "application/json; charset=utf-8", b)
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

func (f *FileKey) Revoke() error {
	message, err := json.Marshal(f)
	signature, _ := sign(ClientPrivateKey, message)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(SignedRequest{message, signature})
	res, err := http.Post(Server+"/revokefile", "application/json; charset=utf-8", b)
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

func GetFileKey(owner string, filename string) (filekey *FileKey, err error) {
	res, err := http.Get(Server + "/users/" + owner + "/" + filename + "/" + ClientUser)
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
	err = json.NewDecoder(res.Body).Decode(&filekey)
	return
}
