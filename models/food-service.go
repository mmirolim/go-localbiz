package models

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type FoodService struct {
	ID         bson.ObjectId `bson:"_id"`
	Address    `bson:"addr"`
	Name       string   `bson:"name"`
	City       string   `bson:"city" json:"city"`
	Desc       string   `bson:"desc"`
	DressCode  string   `bson:"drc"`
	Fax        string   `bson:"fax"`
	Email      string   `bson:"email"`
	OrderPhone string   `bson:"ordtel"`
	WorkHours  string   `bson:"whrs"`
	Halls      string   `bson:"halls"`
	Company    string   `bson:"comp"`
	Cabins     string   `bson:"cabs"`
	Cuisines   []string `bson:"cuisines"`
	Sits       int16    `bson:"sits"`
	Music      []string `bson:"music"`
	Features   []string `bson:"features"`
	Parking    string   `bson:"park"`
	Site       string   `bson:"site"`
	Phones     []string `bson:"tels"`
	Terminal   string   `bson:"term"`
	Types      []string `bson:"types"`
	Transport  string   `bson:"tran"`
	GoodFor    []string `bson:"good_for"`
	Price      string   `bson:"price"`
	Lang       string   `bson:"lang"`
	Geo        `bson:"loc,omitempty"`
	Slug       string    `bson:"slug"`
	Deleted    bool      `bson:"del"`
	UpdatedAt  time.Time `bson:"up_at"`
	CreatedAt  time.Time `bson:"cr_at"`
	CreatedBy  string    `bson:"cr_by"`
	UpdatedBy  string    `bson:"up_by"`
}

func (f FoodService) GetC() string {
	return "food_services"
}

func (f FoodService) GetIndex() []mgo.Index {
	B := Dic.Bson(&f)
	AB := Dic.Bson(Address{})
	return []mgo.Index{
		mgo.Index{
			Key: []string{B("Name")},
		},
		mgo.Index{
			Key: []string{B("Lang"), "-" + B("Name")},
		},
		mgo.Index{
			Key: []string{B("City"), B("Lang")},
		},
		mgo.Index{
			Key: []string{B("Slug")},
		},
		mgo.Index{
			Key:    []string{B("Slug"), B("Lang")},
			Unique: true,
		},
		mgo.Index{
			Key: []string{"$2dsphere:loc"},
		},
		mgo.Index{
			Key: []string{B("Price"), B("Name")},
		},
		mgo.Index{
			Key: []string{B("GoodFor"), B("Name")},
		},
		mgo.Index{
			Key: []string{B("Music"), B("Name")},
		},
		mgo.Index{
			Key: []string{B("Features"), B("Name")},
		},
		mgo.Index{
			Key: []string{B("Types"), B("Name")},
		},
		mgo.Index{
			Key: []string{B("Address") + "." + AB("City"), B("Name")},
		},
		mgo.Index{
			Key: []string{B("Address") + "." + AB("District"), B("Name")},
		},
		mgo.Index{
			Key: []string{B("Deleted")},
		},
		mgo.Index{
			Key: []string{B("UpdatedBy")},
		},
		mgo.Index{
			Key: []string{B("UpdatedAt")},
		},
		mgo.Index{
			Key: []string{B("CreatedAt")},
		},
		mgo.Index{
			Key: []string{B("CreatedBy")},
		},
	}
}

func (f FoodService) GetLocation() Geo {
	return f.Geo
}
func (f *FoodService) FmtFields() {

	f.City = strings.ToLower(f.City)
}
func (f *FoodService) SetDefaults() {
	if f.CreatedAt.Year() == 1 {
		f.CreatedAt = time.Now()
	}
	if f.UpdatedAt.Year() == 1 {
		f.UpdatedAt = time.Now()
	}
}

func (f *FoodService) Validate(s string, bs bson.M) VErrors {
	var vErrors VErrors
	_ = bs
	_ = s
	return vErrors

}

// proxy to Dic
func (fs *FoodService) Bson(f string) string {
	return Dic.Bson(fs)(f)
}
