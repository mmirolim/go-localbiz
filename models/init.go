package models

import (
	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
)

// connection
var (
	MongoHost   string
	MongoDbName string
	MgoSession  *mgo.Session
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
	MongoHost = beego.AppConfig.String("db::mongohost")
	MongoDbName = beego.AppConfig.String("db::mongodbname")
	session, err := mgo.Dial(MongoHost)
	if err != nil {
		beego.Error("Could not connect to mongo instance", err)
	}
	// import to pkg scope
	MgoSession = session

	// init indexes of models
	var fds FoodService
	_, _ = fds.InitIndex()

}
