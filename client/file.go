package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// File Struct
type File struct {
	Id    string
	Owner string
	Name  string
	Data  []byte
}

// Create New File
func NewFile(owner string, name string, data []byte) *File {
	f := new(File)
	f.Owner = owner
	f.Name = name
	f.Data = data
	return f
}

// Upload file to server
func (f *File) Upload() error {
	message, err := json.Marshal(f)
	// Sign the request
	signature, _ := sign(ClientPrivateKey, message)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(SignedRequest{message, signature})
	res, err := http.Post(Server+"/uploadfile", "application/json; charset=utf-8", b)
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

// Get file from server
func GetFile(owner string, filename string) (file *File, err error) {
	res, err := http.Get(Server + "/users/" + owner + "/" + filename)
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
	err = json.NewDecoder(res.Body).Decode(&file)
	return
}

// Get list of users who have access to file from server
func GetFileUsers(owner string, filename string) (users []string, err error) {
	res, err := http.Get(Server + "/users/" + owner + "/" + filename + "/users")
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
	err = json.NewDecoder(res.Body).Decode(&users)
	return
}