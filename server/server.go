package main

import (
	"log"
	"net/http"

	r "github.com/dancannon/gorethink"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	ren "github.com/unrolled/render"
)

var dbSession *r.Session
var render *ren.Render = ren.New(ren.Options{StreamingJSON: true})
var DBHost, Port string

func init() {
	var err error
	dbSession, err = r.Connect(r.ConnectOpts{
		Address:  DBHost + ":28015",
		Database: "Lab2",
		MaxIdle:  10,
		MaxOpen:  10,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	viper.SetDefault("DBHost", "127.0.0.1")
	viper.SetDefault("Port", "8080")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err.Error())
	}
	DBHost = viper.GetString("DBHost")
	Port = viper.GetString("Port")
}

func main() {
	router := httprouter.New()
	router.POST("/register", register)
	router.POST("/uploadfile", uploadFile)
	router.POST("/sharefile", shareFile)
	router.POST("/revokefile", revokeFile)
	router.GET("/users/:username", getUser)
	router.GET("/users/:username/:filename", getFile)
	router.GET("/users/:username/:filename/users", getFileUsers)
	router.GET("/users/:username/:filename/key/:user", getFileKey)

	server := http.Server{
		Addr:    ":" + Port,
		Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error: %v", err)
	}
}
