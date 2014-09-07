package models

import (
	"time"
	"github.com/astaxie/beego"
	"gopkg.in/mgo.v2/bson"
)

var (
	FoodServices FoodService
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
	Id	string
	Address     		 `bson:"address"`
	Name        string   `bson:"name"`
	Description string   `bson:"description"`
	Dresscode   string   `bson:"dresscode"`
	Fax         string   `bson:"fax"`
	Email       string   `bson:"email"`
	Orderphone  string   `bson:"orderPhone"`
	Workdays    string   `bson:"workHours"`
	Halls       string   `bson:"halls"`
	Company     string   `bson:"company"`
	Cabins      string   `bson:"cabins"`
	Cuisine     []string `bson:"cuisine"`
	Sits        int16    `bson:"sits"`
	Music       []string `bson:"music"`
	Refloc      string   `bson:"refLoc"`
	Features    []string `bson:"features"`
	Parking     string   `bson:"parking"`
	Site        string   `bson:"site"`
	Phones      []string `bson:"tels"`
	Terminal    string   `bson:"terminal"`
	Restype     []string `bson:"types"`
	Transport   string   `bson:"trasport"`
	Goodfor     []string `bson:"goodFor"`
	Price       string   `bson:"price"`
	Lang        string   `bson:"lang"`
	GeoJson     		 `bson:"loc,omitempty"`
	Slug		string   `bson:"slug"`
	Deleted		bool 	 `bson:"deleted"`
	Updated		time.Time`bson:"updatedAt"`
	Created		time.Time`bson:"createdAt"`
	CreatedBy	string   `bson:"createdBy"`
	UpdatedBy	string   `bson:"updatedBy"`
}

func (f FoodService) FindOne(b bson.M) (FoodService, error) {
	var foodService FoodService
	session := Session.Copy()
	defer session.Close()

	foodServices := session.DB(Db).C("foodServices")
	err := foodServices.Find(b).One(&foodService)
	if err != nil {
		beego.Error(err)
	}
	return foodService, err
}
