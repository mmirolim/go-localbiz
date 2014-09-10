package models

import (
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
			Key: []string{"good_for", "name"},
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
			Key: []string{"updated_by"},
		},
		mgo.Index{
			Key: []string{"updated_at"},
		},
		mgo.Index{
			Key: []string{"created_at"},
		},
		mgo.Index{
			Key: []string{"created_by"},
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
	Desc		string   `bson:"desc"`
	DressCode   string   `bson:"dress_code"`
	Fax         string   `bson:"fax"`
	Email       string   `bson:"email"`
	OrderPhone  string   `bson:"order_phone"`
	WorkHours   string   `bson:"work_hours"`
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
	GoodFor     []string `bson:"good_for"`
	Price       string   `bson:"price"`
	Lang        string   `bson:"lang"`
	GeoJson     `bson:"loc,omitempty"`
	Slug        string    `bson:"slug"`
	Deleted     bool      `bson:"deleted"`
	UpdatedAt   time.Time `bson:"updated_at"`
	CreatedAt   time.Time `bson:"created_at"`
	CreatedBy   string    `bson:"created_by"`
	UpdatedBy   string    `bson:"updated_by"`
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

func (f FoodService) GetC() string {
	return "food_services"
}

func (f FoodService) InitIndex() (bool, error) {
	var err error
	sess := Session.Copy()
	defer sess.Close()
	for _, v := range indexes {
		err = sess.DB(MongoDbName).C(f.GetC()).EnsureIndex(v)
		if check("FoodService InitIndex -> ", err) {
			return false, err
		}
	}

	return true, err
}

func (f FoodService) Find(q bson.D) ([]FoodService, error) {
	var fds []FoodService
	session := Session.Copy()
	defer session.Close()

	foodServices := session.DB(MongoDbName).C(f.GetC())
	// limit is important when all used, may consume all memory
	// @todo maybe memory consumption reduces if not all fields retrieved?
	iter := foodServices.Find(q).Limit(5000).Iter()
	err := iter.All(&fds)
	check("FoodService FindOne -> ", err)

	return fds, err
}

func (f FoodService) FindOne(q bson.M) (FoodService, error) {
	var fds FoodService
	session := Session.Copy()
	defer session.Close()

	foodServices := session.DB(MongoDbName).C(f.GetC())
	err := foodServices.Find(q).One(&fds)
	check("FoodService FindOne -> ", err)

	return fds, err
}

// @todo refactor maybe loc.coor should be passed by f?
func (f FoodService) FindNear(min, max int, loc GeoJson) (Near, error) {
	var nfs Near
	session := Session.Copy()
	defer session.Close()

	err := session.DB(MongoDbName).Run(bson.D{
		{"geoNear", f.GetC()},
		{"near", bson.D{{"type", "Point"}, {"coordinates", loc.Coordinates}}},
		{"spherical", true},
		{"minDistance", min},
		{"maxDistance", max},
	}, &nfs)
	check("FoodService FindNear -> ", err)

	return nfs, err
}
