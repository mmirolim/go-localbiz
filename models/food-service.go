package models

import (
	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	FoodServices FoodService
	// define indexes
	indexes = []mgo.Index{
		mgo.Index{
			Key: []string{"name"},
		},
		mgo.Index{
			Key: []string{"lang", "-name"},
		},
		mgo.Index{
			Key: []string{"slug"},
		},
		mgo.Index{
			Key: []string{"slug", "lang"},
			Unique: true,
		},
		mgo.Index{
			Key:  []string{"$2dsphere:loc"},
		},
		mgo.Index{
			Key: []string{"price", "name"},
		},
		mgo.Index{
			Key: []string{"goodFor", "name"},
		},
		mgo.Index{
			Key: []string{"music", "name"},
		},
		mgo.Index{
			Key: []string{"features", "name"},
		},
		mgo.Index{
			Key: []string{"types", "name"},
		},
		mgo.Index{
			Key: []string{"address.city", "name"},
		},
		mgo.Index{
			Key: []string{"address.district", "name"},
		},
		mgo.Index{
			Key: []string{"deleted"},
		},
		mgo.Index{
			Key: []string{"updatedBy"},
		},
		mgo.Index{
			Key: []string{"updatedAt"},
		},
		mgo.Index{
			Key: []string{"createdAt"},
		},
		mgo.Index{
			Key: []string{"createdBy"},
		},
	}
)

// define model structs
type GeoJson struct {
	Type        string    `bson:"type"`
	Coordinates []float32 `bson:"coordinates"`
}

type Address struct {
	City     string `bson:"city"`
	District string `bson:"district"`
	Street   string `bson:"street"`
}

type FoodService struct {
	Id          bson.ObjectId `bson:"_id"`
	Address     `bson:"address"`
	Name        string   `bson:"name"`
	Description string   `bson:"description"`
	DressCode   string   `bson:"dressCode"`
	Fax         string   `bson:"fax"`
	Email       string   `bson:"email"`
	OrderPhone  string   `bson:"orderPhone"`
	WorkHours   string   `bson:"workHours"`
	Halls       string   `bson:"halls"`
	Company     string   `bson:"company"`
	Cabins      string   `bson:"cabins"`
	Cuisines    []string `bson:"cuisines"`
	Sits        int16    `bson:"sits"`
	Music       []string `bson:"music"`
	RefLoc      string   `bson:"refLoc"`
	Features    []string `bson:"features"`
	Parking     string   `bson:"parking"`
	Site        string   `bson:"site"`
	Phones      []string `bson:"tels"`
	Terminal    string   `bson:"terminal"`
	Types       []string `bson:"types"`
	Transport   string   `bson:"trasport"`
	GoodFor     []string `bson:"goodFor"`
	Price       string   `bson:"price"`
	Lang        string   `bson:"lang"`
	GeoJson     `bson:"loc,omitempty"`
	Slug        string    `bson:"slug"`
	Deleted     bool      `bson:"deleted"`
	UpdatedAt   time.Time `bson:"updatedAt"`
	CreatedAt   time.Time `bson:"createdAt"`
	CreatedBy   string    `bson:"createdBy"`
	UpdatedBy   string    `bson:"updatedBy"`
}

// struct to store Near FoodServices result from mongo
type Near struct {
	Results []struct {
		Dis float32
		Obj FoodService
	}
	Stats NearStats
	Ok    float32
}

func check(s string, e error) bool {
	if e != nil {
		beego.Error(s + e.Error())
		return true
	}
	return false
}

func (f FoodService) GetC() string {
	return "foodServices"
}

func (f FoodService) InitIndex() (bool, error) {
	var err error
	sess := Session.Copy()
	defer sess.Close()
	for _, v := range indexes {
		err = sess.DB(Db).C(f.GetC()).EnsureIndex(v)
		if check("FoodService InitIndex -> ", err) {
			return false, err
		}
	}

	return true, err
}

func (f FoodService) FindOne(b bson.M) (FoodService, error) {
	var fds FoodService
	session := Session.Copy()
	defer session.Close()

	foodServices := session.DB(Db).C(fds.GetC())
	err := foodServices.Find(b).One(&fds)
	check("FoodService FindOne -> ", err)

	return fds, err
}

func (f FoodService) FindNear(min, max int, loc GeoJson) (Near, error) {
	var nfs Near
	session := Session.Copy()
	defer session.Close()

	err := session.DB(Db).Run(bson.D{
		{"geoNear", "foodServices"},
		{"near", bson.D{{"type", "Point"}, {"coordinates", loc.Coordinates}}},
		{"spherical", true},
		{"minDistance", min},
		{"maxDistance", max},
	}, &nfs)
	check("FoodService FindNear -> ", err)

	return nfs, err
}
