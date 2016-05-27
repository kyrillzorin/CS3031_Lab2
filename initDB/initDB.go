package main

import (
	"log"

	r "github.com/dancannon/gorethink"
	"github.com/spf13/viper"
)

var dbSession *r.Session
var DBHost string

func init() {
	var err error
	dbSession, err = r.Connect(r.ConnectOpts{
		Address: DBHost + ":28015",
		MaxIdle: 10,
		MaxOpen: 10,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	viper.SetDefault("DBHost", "127.0.0.1")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}
	DBHost = viper.GetString("DBHost")
}

func main() {
	_, err := r.DBCreate("Lab2").RunWrite(dbSession)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = r.DB("Lab2").TableCreate("users").RunWrite(dbSession)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = r.DB("Lab2").Table("users").IndexCreate("username").RunWrite(dbSession)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = r.DB("Lab2").TableCreate("files").RunWrite(dbSession)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = r.DB("Lab2").Table("files").IndexCreate("name").RunWrite(dbSession)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = r.DB("Lab2").Table("files").IndexCreate("owner").RunWrite(dbSession)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = r.DB("Lab2").TableCreate("filekeys").RunWrite(dbSession)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = r.DB("Lab2").Table("filekeys").IndexCreate("name").RunWrite(dbSession)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = r.DB("Lab2").Table("filekeys").IndexCreate("owner").RunWrite(dbSession)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = r.DB("Lab2").Table("filekeys").IndexCreate("user").RunWrite(dbSession)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
