package main

import (
	"crypto/rsa"
	"errors"

	r "github.com/dancannon/gorethink"
)

// User DB table
var userTable r.Term = r.Table("users")

type User struct {
	Id       string         `gorethink:"id,omitempty"`
	Username string         `gorethink:"username"`
	PubKey   *rsa.PublicKey `gorethink:"pubkey"`
}

// Inserts user into DB
func (u *User) Insert(dbSession *r.Session) (wRes r.WriteResponse, err error) {
	res, err := userTable.GetAllByIndex("username", u.Username).Run(dbSession)
	if err != nil {
		return wRes, err
	}
	if !res.IsNil() {
		return wRes, errors.New("Duplicate account")
	}
	wRes, err = userTable.Insert(u).RunWrite(dbSession)
	return
}

// Gets a user from the DB
func GetUser(username string, dbSession *r.Session) (user *User, err error) {
	res, err := userTable.GetAllByIndex("username", username).Run(dbSession)
	if err != nil {
		return
	}
	user = new(User)
	err = res.One(&user)
	return
}
