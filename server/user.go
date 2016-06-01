package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"errors"

	r "github.com/dancannon/gorethink"
)

// User DB table
var userTable r.Term = r.Table("users")

type User struct {
	Id       string
	Username string
	PubKey   *rsa.PublicKey
}

type dbUser struct {
	Id       string `gorethink:"id,omitempty"`
	Username string `gorethink:"username"`
	PubKey   []byte `gorethink:"pubkey"`
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
	var user dbUser
	user.Id = u.Id
	user.Username = u.Username
	var pubKey bytes.Buffer
	enc := gob.NewEncoder(&pubKey)
	err = enc.Encode(u.PubKey)
	if err != nil {
		return wRes, err
	}
	user.PubKey = pubKey.Bytes()
	wRes, err = userTable.Insert(user).RunWrite(dbSession)
	return
}

// Gets a user from the DB
func GetUser(username string, dbSession *r.Session) (user *User, err error) {
	res, err := userTable.GetAllByIndex("username", username).Run(dbSession)
	if err != nil {
		return
	}
	u := new(dbUser)
	err = res.One(&u)
	if err != nil {
		return
	}
	user = new(User)
	user.Id = u.Id
	user.Username = u.Username
	pubKey := bytes.NewBuffer(u.PubKey)
	dec := gob.NewDecoder(pubKey)
	err = dec.Decode(&user.PubKey)
	return
}
