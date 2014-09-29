package models

import (
	"github.com/astaxie/beego/validation"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

// roles admin = 1, editor = 2, tester = 3, client = 4, user = 5
const (
	// roles
	_ = iota
	roleAdmin
	roleEditor
	roleTester
	roleClient
	roleUser
)

var (
	// define indexes
	userIndexes = []mgo.Index{
		mgo.Index{
			Key:    []string{"username"},
			Unique: true,
		},
		mgo.Index{
			Key: []string{"name"},
		},
		mgo.Index{
			Key:    []string{"fb_data.id"},
			Unique: true,
		},
		mgo.Index{
			Key: []string{"gender"},
		},

		mgo.Index{
			Key: []string{"role"},
		},
		mgo.Index{
			Key: []string{"b_day"},
		},

		mgo.Index{
			Key:    []string{"email"},
			Unique: true,
		},

		mgo.Index{
			Key: []string{"last_login_at"},
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

//@todo fix validation for email now requiring
type User struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	UserName     string        `bson:"username" json:"username" valid:"Required; MinSize(6); MaxSize(100)"`
	Email        string        `bson:"email" json:"email" valid:"MinSize(4);MaxSize(100)"`
	Name         string        `bson:"name" json:"name" valid:"Required; MaxSize(100)"`
	FirstName    string        `bson:"first_name" json:"first_name" valid:"Required; MaxSize(100)"`
	LastName     string        `bson:"last_name" json:"last_name" valid:"Required; MaxSize(100)"`
	Gender       string        `bson:"gender" json:"gender"  valid:"MaxSize(20)"`
	Locale       string        `bson:"locale" json:"locale" json:"locale"  valid:"Length(5)"`
	LastLoginAt  time.Time     `bson:"last_login_at" json:"last_login_at"`
	Role         int           `bson:"role" json:"role"`
	Bday         time.Time     `bson:"bday,omitempty" json:"bday,omitempty"`
	FacebookData `bson:"fb_data" json:"fb_data"`
	Address      `bson:"address" json:"address"`
	Geo          `bson:"loc,omitempty" json:"loc,omitempty`
	Deleted      bool      `bson:"deleted" json:"deleted"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	CreatedBy    string    `bson:"created_by" json:"created_by"`
	UpdatedBy    string    `bson:"updated_by" json:"updated_by"`
	IsAdmin      bool      `bson:"is_admin,omitempty" json:"is_admin,omitempty"`
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

	u.UpdatedAt = time.Now()

	if u.CreatedAt.Year() == 1 {
		u.CreatedAt = time.Now()
	}
	if u.LastLoginAt.Year() == 1 {
		u.LastLoginAt = time.Now()
	}
	if u.Role == 0 {
		u.Role = roleUser
	}
}

// validation custom called after tagfield validation
func (u *User) Valid(v *validation.Validation) {
	//@todo maybe make regex
	if strings.Index(u.UserName, "admin") != -1 && u.Role != roleAdmin {
		v.SetError("Username", "Can't contain 'admin' in UserName")
	}

}
