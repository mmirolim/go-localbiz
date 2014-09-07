package models

import (
	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
)

// connection
var (
	Host string
	Db	 string
	Session *mgo.Session
)

func InitConnection() {
	Host = beego.AppConfig.String("db::host")
	Db = beego.AppConfig.String("db::db")
	session, err := mgo.Dial(Host)
	if err != nil {
		beego.Error("Could not connect to mongo instance", err)
	}
	// import to pkg scope
	Session = session
}
