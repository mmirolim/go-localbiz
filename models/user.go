package models

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	// define indexes
	mgoIndexes = []mgo.Index{
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
		Key:    []string{"slug", "lang"},
		Unique: true,
	},
	mgo.Index{
		Key: []string{"$2dsphere:loc"},
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

type FoodService struct {
	Id         bson.ObjectId `bson:"_id"`
	Address    `bson:"address"`
	Name       string   `bson:"name"`
	Desc       string   `bson:"desc"`
	DressCode  string   `bson:"dress_code"`
	Fax        string   `bson:"fax"`
	Email      string   `bson:"email"`
	OrderPhone string   `bson:"order_phone"`
	WorkHours  string   `bson:"work_hours"`
	Halls      string   `bson:"halls"`
	Company    string   `bson:"company"`
	Cabins     string   `bson:"cabins"`
	Cuisines   []string `bson:"cuisines"`
	Sits       int16    `bson:"sits"`
	Music      []string `bson:"music"`
	Features   []string `bson:"features"`
	Parking    string   `bson:"parking"`
	Site       string   `bson:"site"`
	Phones     []string `bson:"tels"`
	Terminal   string   `bson:"terminal"`
	Types      []string `bson:"types"`
	Transport  string   `bson:"trasport"`
	GoodFor    []string `bson:"good_for"`
	Price      string   `bson:"price"`
	Lang       string   `bson:"lang"`
	Geo        `bson:"loc,omitempty"`
	Slug       string    `bson:"slug"`
	Deleted    bool      `bson:"deleted"`
	UpdatedAt  time.Time `bson:"updated_at"`
	CreatedAt  time.Time `bson:"created_at"`
	CreatedBy  string    `bson:"created_by"`
	UpdatedBy  string    `bson:"updated_by"`
}

func (f FoodService) GetC() string {
	return "users"
}

func (f FoodService) GetIndex() []mgo.Index {
	return mgoIndexes
}

func (f FoodService) GetLocation() Geo {
	return f.Geo
}
