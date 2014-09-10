package models

import (
	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
)

// connection
var (
	Host    string
	Db      string
	Session *mgo.Session
)

type NearStats struct {
	NScanned  uint32  `bson:"nscanned"`
	ObjLoaded uint32  `bson:"objectsLoaded"`
	AvrDis    float32 `bson:"avgDistance"`
	MaxDis    float32 `bson:"maxDistance"`
	time      int32   `bson:"time"`
}

func check(s string, e error) bool {
	if e != nil {
		beego.Error(s + e.Error())
		return true
	}
	return false
}

func InitConnection() {
	Host = beego.AppConfig.String("db::host")
	Db = beego.AppConfig.String("db::db")
	session, err := mgo.Dial(Host)
	if err != nil {
		beego.Error("Could not connect to mongo instance", err)
	}
	// import to pkg scope
	Session = session

	// init indexes of models
	var fds FoodService
	_, _ = fds.InitIndex()

}
