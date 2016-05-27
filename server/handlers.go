package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func indexHandler(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	http.Redirect(w, req, "/lol", http.StatusFound)
}

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
	if !verify(user.PubKey, signedRequest.Message, signedRequest.Signature) {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	_, err = file.Insert(dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"Status": "success", "Error": ""})
}

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
	if !verify(user.PubKey, signedRequest.Message, signedRequest.Signature) {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	_, err = filekey.Insert(dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"Status": "success", "Error": ""})
}

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
	if !verify(user.PubKey, signedRequest.Message, signedRequest.Signature) {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	_, err = filekey.Revoke(dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"Status": "success", "Error": ""})
}

func getUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	user, err := GetUser(ps.ByName("username"), dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, user)
}

func getFile(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	file, err := GetFile(ps.ByName("username"), ps.ByName("filename"), dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, file)
}

func getFileUsers(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	users, err := GetFileUsers(ps.ByName("username"), ps.ByName("filename"), dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, users)
}

func getFileKey(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	filekey, err := GetFileKey(ps.ByName("username"), ps.ByName("filename"), ps.ByName("user"), dbSession)
	if err != nil {
		render.JSON(w, http.StatusBadRequest, map[string]string{"Status": "failure", "Error": err.Error()})
		return
	}
	render.JSON(w, http.StatusOK, filekey)
}
