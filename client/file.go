package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// File Struct
type File struct {
	Id    string
	Owner string
	Name  string
	Data  []byte
}

// File Users Struct
type FileUsers struct {
	Users []string
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
	res, err := http.Post(Server+"/uploadfile", "application/json; charset=utf-8", b)
	defer res.Body.Close()
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
	defer res.Body.Close()
	if err != nil {
		return
	}
	if res.Body == nil {
		err = errors.New("Empty Response")
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
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&file)
	return
}

// Get list of users who have access to file from server
func GetFileUsers(owner string, filename string) (users []string, err error) {
	res, err := http.Get(Server + "/users/" + owner + "/" + filename + "/users")
	defer res.Body.Close()
	if err != nil {
		return
	}
	if res.Body == nil {
		err = errors.New("Empty Response")
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
	userList := new(FileUsers)
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&userList)
	if err != nil {
		return
	}
	users = userList.Users
	return
}
