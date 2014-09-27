package models

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	// roles
	roleAdmin = iota
	roleUser
	roleEditor
	roleTester
	roleClient
)

var (
	// define indexes
	userIndexes = []mgo.Index{
		mgo.Index{
			Key: []string{"name"},
		},
		mgo.Index{
			Key: []string{"lang", "-name"},
		},
		mgo.Index{
			Key:    []string{"slug", "lang"},
			Unique: true,
		},
		mgo.Index{
			Key: []string{"$2dsphere:loc"},
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

type FacebookData struct {
	Id          string `bson:"id" json:"id`
	Link        string `bson:"link" json:"link"`
	Name        string `bson:"name" json:"name"`
	FirstName   string `bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name"`
	Gender      string `bson:"gender" json:"gender"`
	UserName    string `bson:"username" json:"username`
	Locale      string `bson:"locale" json:"locale"`
	AccessToken string `bson:"access_token"`
}

type User struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	UserName     string        `bson:"username" json:"username"`
	Email        string        `bson:"email" json:"email"`
	Name         string        `bson:"name" json:"name"`
	FirstName    string        `bson:"first_name" json:"first_name"`
	LastName     string        `bson:"last_name" json:"last_name"`
	Gender       string        `bson:"gender" json:"gender"`
	Locale       string        `bson:"locale" json:"locale" json:"locale"`
	LastLoginAt  time.Time     `bson:"last_login_at" json:"last_login_at"`
	Role         int           `bson:"role" json:"role"`
	FacebookData `bson:"fb_data" json:"fb_data"`
	Address      `bson:"address" json:"address"`
	Geo          `bson:"loc,omitempty" json:"loc,omitempty`
	Deleted      bool      `bson:"deleted" json:"deleted"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	CreatedBy    string    `bson:"created_by" json:"created_by"`
	UpdatedBy    string    `bson:"updated_by" json:"updated_by"`
	IsAdmin      bool      `bson:is_admin,omitempty" json:"is_admin,omitempty"`
}

func (u User) GetC() string {
	return "users"
}

func (u User) GetIndex() []mgo.Index {
	return userIndexes
}

func (u User) GetLocation() Geo {
	return u.Geo
}

func (u *User) SetDefaults() {
	if u.CreatedAt.Year() == 1 {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.Year() == 1 {
		u.UpdatedAt = time.Now()
	}
	if u.Role == 0 {
		u.Role = roleUser
	}
}
