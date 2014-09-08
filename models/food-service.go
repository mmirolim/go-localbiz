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
	DressCode   string   `bson:"dressCode"`
	Fax         string   `bson:"fax"`
	Email       string   `bson:"email"`
	OrderPhone  string   `bson:"orderPhone"`
	WorkHours    string   `bson:"workHours"`
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
	Types        []string `bson:"types"`
	Transport   string   `bson:"trasport"`
	GoodFor     []string `bson:"goodFor"`
	Price       string   `bson:"price"`
	Lang        string   `bson:"lang"`
	GeoJson     		 `bson:"loc,omitempty"`
	Slug		string   `bson:"slug"`
	Deleted		bool 	 `bson:"deleted"`
	UpdatedAt		time.Time `bson:"updatedAt"`
	CreatedAt		time.Time `bson:"createdAt"`
	CreatedBy	string   `bson:"createdBy"`
	UpdatedBy	string   `bson:"updatedBy"`
}
type NearStats struct {
	NScanned uint32 `bson:"nscanned"`
	ObjLoaded uint32 `bson:"objectsLoaded"`
	AvrDis	float32	`bson:"avgDistance"`
	MaxDis	float32	`bson:"maxDistance"`
	time	int32	`bson:"time"`
}
// struct to store Near FoodServices result from mongo
type NearFoodService struct {
	Dis	float32	`bson:"dis"`
	Obj	FoodService `bson:"obj"`
}

type NearResult struct {
	Results []NearFoodService
	Stats	NearStats
	Ok	float32
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

func (f FoodService) FindNear(min, max int, loc GeoJson) (NearResult, error) {
	var nfs NearResult
	session := Session.Copy()
	defer session.Close()

	err := session.DB(Db).Run(bson.D{
		{"geoNear", "foodServices" },
		{"near", bson.D{ {"type", "Point"}, {"coordinates", loc.Coordinates}}},
		{"spherical", true},
		{"minDistance", min},
		{"maxDistance", max},
		}, &nfs)
	if err != nil {
		beego.Error(err)
	}

	return nfs, err
}

