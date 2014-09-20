package models

import (
	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// connection
var (
	MongoHost   string
	MongoDbName string
	MgoSession  *mgo.Session
)

type BsonData struct {
	Raw bson.Raw
}

type DocModel interface {
	GetC() string
	GetIndex() []mgo.Index
	GetLocation() Geo
}

// define model structs
type Geo struct {
	Type        string    `bson:"type"`
	Coordinates []float32 `bson:"coordinates"`
}

type Address struct {
	City     string `bson:"city"`
	District string `bson:"district"`
	Street   string `bson:"street"`
	RefLoc   string `bson:"ref_loc"`
}

type NearStats struct {
	NScanned  uint32  `bson:"nscanned"`
	ObjLoaded uint32  `bson:"objectsLoaded"`
	AvrDis    float32 `bson:"avgDistance"`
	MaxDis    float32 `bson:"maxDistance"`
	time      int32   `bson:"time"`
}

// struct to store Near FoodServices result from mongo
type Near struct {
	Results bson.Raw
	Stats   NearStats
	Ok      float32
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
	_, _ = DocInitIndex(fds)

}

func DocInitIndex(m DocModel) (bool, error) {
	var err error
	sess := MgoSession.Copy()
	defer sess.Close()
	for _, v := range m.GetIndex() {
		err = sess.DB(MongoDbName).C(m.GetC()).EnsureIndex(v)
		if check("InitIndex -> ", err) {
			return false, err
		}
	}

	return true, err
}

func DocFind(q bson.D, m DocModel, data interface{}) error {
	session := MgoSession.Copy()
	defer session.Close()

	foodServices := session.DB(MongoDbName).C(m.GetC())
	// limit is important when all used, may consume all memory
	// @todo maybe memory consumption reduces if not all fields retrieved?
	iter := foodServices.Find(q).Limit(5000).Iter()
	err := iter.All(data)
	check("FindOne -> ", err)

	return err
}

func DocFindOne(q bson.M, m DocModel) error {
	session := MgoSession.Copy()
	defer session.Close()

	foodServices := session.DB(MongoDbName).C(m.GetC())
	err := foodServices.Find(q).One(m)
	check("FindOne -> ", err)

	return err
}

// @todo refactor maybe loc.coor should be passed by f?
func DocFindNear(min, max int, m DocModel, data interface{}) error {
	session := MgoSession.Copy()
	defer session.Close()

	err := session.DB(MongoDbName).Run(bson.D{
		{"geoNear", m.GetC()},
		{"near", bson.D{{"type", "Point"}, {"coordinates", m.GetLocation().Coordinates}}},
		{"spherical", true},
		{"minDistance", min},
		{"maxDistance", max},
	}, data)
	check("FindNear -> ", err)

	return err
}
