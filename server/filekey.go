package main

import (
	"errors"

	r "github.com/dancannon/gorethink"
)

// File Key DB table
var fileKeyTable r.Term = r.Table("filekeys")

// File Key Struct
type FileKey struct {
	Id    string `gorethink:"id,omitempty"`
	User  string `gorethink:"user"`
	Owner string `gorethink:"owner"`
	Name  string `gorethink:"name"`
	Key   []byte `gorethink:"key"`
}

// File Users Struct
type FileUsers struct {
	Users []string
}

// Inserts file key into DB
func (f *FileKey) Insert(dbSession *r.Session) (res r.WriteResponse, err error) {
	dbRes, err := fileKeyTable.GetAllByIndex("name", f.Name).Filter(map[string]interface{}{"owner": f.Owner, "user": f.User}).Run(dbSession)
	if err != nil {
		return
	}
	if !dbRes.IsNil() {
		filekey := new(FileKey)
		err = dbRes.One(&filekey)
		if err != nil {
			return
		}
		f.Id = filekey.Id
		res, err = f.Update(dbSession)
		return
	}
	res, err = fileKeyTable.Insert(f).RunWrite(dbSession)
	return
}

// Updates file key in DB
func (f *FileKey) Update(dbSession *r.Session) (res r.WriteResponse, err error) {
	res, err = fileKeyTable.Get(f.Id).Update(f).RunWrite(dbSession)
	return
}

// Revoke file key from DB
func (f *FileKey) Revoke(dbSession *r.Session) (res r.WriteResponse, err error) {
	if f.User == f.Owner {
		err = errors.New("Can't revoke own file access")
		return
	}
	dbRes, err := fileKeyTable.GetAllByIndex("name", f.Name).Filter(map[string]interface{}{"owner": f.Owner, "user": f.User}).Run(dbSession)
	if err != nil {
		return
	}
	if !dbRes.IsNil() {
		filekey := new(FileKey)
		err = dbRes.One(&filekey)
		if err != nil {
			return
		}
		res, err = fileKeyTable.Get(filekey.Id).Delete().RunWrite(dbSession)
		return
	}
	return
}

// Get file key from DB
func GetFileKey(owner string, filename string, user string, dbSession *r.Session) (filekey *FileKey, err error) {
	res, err := fileKeyTable.GetAllByIndex("name", filename).Filter(map[string]interface{}{"owner": owner, "user": user}).Run(dbSession)
	if err != nil {
		return
	}
	if res.IsNil() {
		err = errors.New("You do not have access to this file")
		return
	}
	filekey = new(FileKey)
	err = res.One(&filekey)
	return
}

// Get a slice (array) of users who have keys to the file
func GetFileUsers(owner string, filename string, dbSession *r.Session) (userList *FileUsers, err error) {
	res, err := fileKeyTable.GetAllByIndex("name", filename).Filter(map[string]interface{}{"owner": owner}).Pluck("user").Run(dbSession)
	if err != nil {
		return
	}
	var userMap []map[string]string
	err = res.All(&userMap)
	if err != nil {
		return
	}
	users := make([]string, 0, len(userMap))
	for _, user := range userMap {
		users = append(users, user["user"])
	}
	userList = new(FileUsers)
	userList.Users = users
	return
}
