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

// Inserts file key into DB
func (f *FileKey) Insert(dbSession *r.Session) (res r.WriteResponse, err error) {
	dbRes, err := fileKeyTable.GetAllByIndex("name", f.Name).GetAllByIndex("owner", f.Owner).GetAllByIndex("user", f.User).Run(dbSession)
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
	dbRes, err := fileKeyTable.GetAllByIndex("name", f.Name).GetAllByIndex("owner", f.Owner).GetAllByIndex("user", f.User).Run(dbSession)
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
	res, err := fileKeyTable.GetAllByIndex("name", filename).GetAllByIndex("owner", owner).GetAllByIndex("user", user).Run(dbSession)
	if err != nil {
		return
	}
	filekey = new(FileKey)
	err = res.One(&filekey)
	return
}

// Get a slice (array) of users who have keys to the file
func GetFileUsers(owner string, filename string, dbSession *r.Session) (users []string, err error) {
	res, err := fileKeyTable.GetAllByIndex("name", filename).GetAllByIndex("owner", owner).Pluck("user").Run(dbSession)
	if err != nil {
		return
	}
	err = res.All(&users)
	return
}
