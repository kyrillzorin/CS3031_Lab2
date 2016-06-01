package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Register a new user
func register(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var user User
	if req.Body == nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": "Invalid Request: Empty"})
		return
	}
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	_, err = user.Insert(dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"Status": "success", "Error": ""})
}

// Handle a file upload
func uploadFile(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var signedRequest SignedRequest
	if req.Body == nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": "Invalid Request: Empty"})
		return
	}
	err := json.NewDecoder(req.Body).Decode(&signedRequest)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	var file File
	err = json.Unmarshal(signedRequest.Message, &file)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	user, _ := GetUser(file.Owner, dbSession)
	// Verify signed message
	if !verify(user.PubKey, signedRequest.Message, signedRequest.Signature) {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": "Could not verify signature"})
		return
	}
	_, err = file.Insert(dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"Status": "success", "Error": ""})
}

// Share file access with a user
func shareFile(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var signedRequest SignedRequest
	if req.Body == nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": "Invalid Request: Empty"})
		return
	}
	err := json.NewDecoder(req.Body).Decode(&signedRequest)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	var filekey FileKey
	err = json.Unmarshal(signedRequest.Message, &filekey)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	user, _ := GetUser(filekey.Owner, dbSession)
	// Verify signed message
	if !verify(user.PubKey, signedRequest.Message, signedRequest.Signature) {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": "Could not verify signature"})
		return
	}
	_, err = filekey.Insert(dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"Status": "success", "Error": ""})
}

// Revoke file access for a user
func revokeFile(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var signedRequest SignedRequest
	if req.Body == nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": "Invalid Request: Empty"})
		return
	}
	err := json.NewDecoder(req.Body).Decode(&signedRequest)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	var filekey FileKey
	err = json.Unmarshal(signedRequest.Message, &filekey)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	user, _ := GetUser(filekey.Owner, dbSession)
	// Verify signed message
	if !verify(user.PubKey, signedRequest.Message, signedRequest.Signature) {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": "Could not verify signature"})
		return
	}
	_, err = filekey.Revoke(dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"Status": "success", "Error": ""})
}

// Get a user
func getUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	user, err := GetUser(ps.ByName("username"), dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, user)
}

// Get a file
func getFile(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	file, err := GetFile(ps.ByName("username"), ps.ByName("filename"), dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, file)
}

// Get a list of users with access to a file
func getFileUsers(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	users, err := GetFileUsers(ps.ByName("username"), ps.ByName("filename"), dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, users)
}

// Get a file key
func getFileKey(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	filekey, err := GetFileKey(ps.ByName("username"), ps.ByName("filename"), ps.ByName("user"), dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, filekey)
}
