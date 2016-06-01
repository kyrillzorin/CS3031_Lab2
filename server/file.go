package main

import (
	"errors"

	r "github.com/dancannon/gorethink"
)

// File DB table
var fileTable r.Term = r.Table("files")

// File Struct
type File struct {
	Id    string `gorethink:"id,omitempty"`
	Owner string `gorethink:"owner"`
	Name  string `gorethink:"name"`
	Data  []byte `gorethink:"data"`
}

// Inserts file into DB, Updates file if it already exists
func (f *File) Insert(dbSession *r.Session) (res r.WriteResponse, err error) {
	dbRes, err := fileTable.GetAllByIndex("name", f.Name).Filter(map[string]interface{}{"owner": f.Owner}).Run(dbSession)
	if err != nil {
		return
	}
	if !dbRes.IsNil() {
		file := new(File)
		err = dbRes.One(&file)
		if err != nil {
			return
		}
		f.Id = file.Id
		res, err = f.Update(dbSession)
		return
	}
	res, err = fileTable.Insert(f).RunWrite(dbSession)
	return
}

// Updates file in DB
func (f *File) Update(dbSession *r.Session) (res r.WriteResponse, err error) {
	res, err = fileTable.Get(f.Id).Update(f).RunWrite(dbSession)
	return
}

// Get a file from DB
func GetFile(owner string, filename string, dbSession *r.Session) (file *File, err error) {
	res, err := fileTable.GetAllByIndex("name", filename).Filter(map[string]interface{}{"owner": owner}).Run(dbSession)
	if err != nil {
		return
	}
	if res.IsNil() {
		err = errors.New("File does not exist")
		return
	}
	file = new(File)
	err = res.One(&file)
	return
}
