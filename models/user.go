package models

import (
	"github.com/astaxie/beego"
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
	Id          string `bson:"id" json:"id"`
	Link        string `bson:"link" json:"link"`
	Name        string `bson:"name" json:"name"`
	FirstName   string `bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name"`
	Gender      string `bson:"gender" json:"gender"`
	UserName    string `bson:"username" json:"username"`
	Locale      string `bson:"locale" json:"locale"`
	AccessToken string `bson:"access_token"`
}

type GoogleData struct {
	Id         string `bson:"id" json:"id"`
	ObjectType string `bson:"objectType" json:"objectType"`
	Kind       string `bson:"kind" json:"kind"`
	Etag       string `bson:"etag" json:"etag"`
	PlaceLived struct {
		Value   string `bson:"value" json:"value"`
		Primary bool   `bson:"primary" json:"primary"`
	}
	DisplayName string `bson:"displayName" json:"displayName"`
	Url         string `bson:"url" json:"url"`
	Name        struct {
		FamilyName string `bson:"familyName" json:"familyName"`
		GivenName  string `bson:"givenName" json:"givenName"`
	} `bson:"name" json:"name"`
	Image struct {
		Url       string `bson:"url" json:"url"`
		IsDefault bool   `bson:"isDefault" json:"isDefault"`
	} `bson:"image" json:"image"`
	FirstName   string `bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name"`
	Gender      string `bson:"gender" json:"gender"`
	UserName    string `bson:"username" json:"username"`
	Lang        string `bson:"language" json:"language"`
	isPlusUser  string `bosn:"is_plus_user" json:"isPlusUser"`
	Verified    bool   `bson:"verified" json:"verified"`
	AccessToken string `bson:"access_token"`
}

// @important  if Field naming changes validation also should be changed
//@todo fix validation for email now requiring
type User struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	UserName     string        `bson:"username" json:"username"`
	Email        string        `bson:"email" json:"email"`
	Name         string        `bson:"name" json:"name"`
	City         string        `bson:"city" json:"city"`
	FirstName    string        `bson:"first_name" json:"first_name"`
	LastName     string        `bson:"last_name" json:"last_name"`
	Gender       string        `bson:"gender" json:"gender"`
	Locale       string        `bson:"locale" json:"locale" json:"locale"`
	LastLoginAt  time.Time     `bson:"last_login_at" json:"last_login_at"`
	Role         int           `bson:"role" json:"role"`
	Bday         time.Time     `bson:"bday,omitempty" json:"bday,omitempty"`
	FacebookData `bson:"fb_data" json:"fb_data"`
	GoogleData   `bson:"gg_data" json:"gg_data"`
	Address      `bson:"address" json:"address"`
	Geo          `bson:"loc,omitempty" json:"loc,omitempty"`
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

func (u *User) FmtFields() {
	//@todo create helper fmt function to pass func name to run on argument
	actions := []string{"ToLower", "TrimSpace"}
	trimSpace := []string{"TrimSpace"}
	u.UserName = FmtString(u.UserName, actions)
	u.FirstName = FmtString(u.FirstName, trimSpace)
	u.LastName = FmtString(u.LastName, trimSpace)
	u.Email = FmtString(u.Email, actions)
	u.Locale = FmtString(u.Locale, actions)
	u.City = FmtString(u.City, actions)

}

func (u *User) InitWithFb(fb FacebookData) {
	u.FacebookData = fb
	u.UserName = fb.UserName
	u.FirstName = fb.FirstName
	u.LastName = fb.LastName
	u.Locale = strings.ToLower(fb.Locale)
	u.Name = fb.Name
	u.Gender = fb.Gender
}

func (u *User) InitWithGg(gg GoogleData) {
	u.GoogleData = gg
	u.FirstName = gg.Name.GivenName
	u.LastName = gg.Name.FamilyName
	u.SetName(u.FirstName, u.LastName)
	u.Gender = gg.Gender
	if gg.PlaceLived.Value != "" {
		u.City = strings.ToLower(gg.PlaceLived.Value)
	}
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

func (u *User) SetName(firstName, lastName string) {
	u.Name = strings.TrimSpace(firstName) + " " + strings.TrimSpace(lastName)
}

// return bson field name from cached FieldDic, convenience func
func (u User) Bson(f string) string {
	b, ok := FieldDic["User"]["FieldBson"][f]
	if !ok {
		beego.Error("User.Bson key in FieldDic does not exists " + f)
	}

	return b
}

func (u User) Field(b string) string {
	f, ok := FieldDic["User"]["BsonField"][b]
	if !ok {
		beego.Error("User.Field key in FieldDic does not exists " + b)
	}
	return f
}

// validate field of DocModel
//@todo all msg should be translatable
func (u *User) Validate(bs bson.M) VErrors {
	v := Validator{}
	// get bson field name
	b := u.Bson
	// if bs empty validatate all fields
	// else validate only updated fields
	l := len(bs)
	f := "UserName"
	// check if field has udpate
	_, ok := bs[b(f)]
	if l == 0 || ok {
		v.Required(u.UserName, b(f))
		v.Size(u.UserName, b(f), 2, 100)
		v.AlphaDash(u.UserName, b(f))
		v.NotContainStr(u.UserName, b(f), []string{"admin", "administrator", "админ", "администратор"})
	}

	f = "LastName"
	_, ok = bs[b(f)]
	if l == 0 || ok {
		v.Required(u.LastName, b(f))
		v.Size(u.LastName, b(f), 2, 100)
	}

	f = "FirstName"
	_, ok = bs[b(f)]
	if l == 0 || ok {
		v.Required(u.FirstName, b(f))
		v.Size(u.FirstName, b(f), 2, 100)
	}
	f = "Email"
	_, ok = bs[b(f)]
	// not required field should validated if not empty
	if u.Email != "" || ok {
		v.Email(u.Email, b(f))
		v.Size(u.Email, b(f), 5, 100)
	}
	return v.Errors
}
