package models

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/utils"
	"github.com/nicksnyder/go-i18n/i18n"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"strings"
	"time"
)

// roles admin = 1, editor = 2, tester = 3, client = 4, user = 5
const (

	roleAdmin = 1
	roleEditor = 2
	roleTester = 3
	roleClient = 4
	roleUser = 5
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

type FBData struct {
	ID          string `bson:"id" json:"id"`
	Link        string `bson:"link" json:"link"`
	Name        string `bson:"name" json:"name"`
	FirstName   string `bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name"`
	Gender      string `bson:"gender" json:"gender"`
	UserName    string `bson:"username" json:"username"`
	Locale      string `bson:"locale" json:"locale"`
	AccessToken string `bson:"access_token"`
}

type GGData struct {
	ID         string `bson:"id" json:"id"`
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
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	UserName    string        `bson:"username" json:"username"`
	Email       string        `bson:"email" json:"email"`
	Name        string        `bson:"name" json:"name"`
	City        string        `bson:"city" json:"city"`
	FirstName   string        `bson:"first_name" json:"first_name"`
	LastName    string        `bson:"last_name" json:"last_name"`
	Gender      string        `bson:"gender" json:"gender"`
	Locale      string        `bson:"locale" json:"locale" json:"locale"`
	LastLoginAt time.Time     `bson:"last_login_at" json:"last_login_at"`
	Role        int           `bson:"role" json:"role"`
	Bday        time.Time     `bson:"bday,omitempty" json:"bday,omitempty"`
	FBData      `bson:"fb_data" json:"fb_data"`
	GGData      `bson:"gg_data" json:"gg_data"`
	Address     `bson:"address" json:"address"`
	Geo         `bson:"loc,omitempty" json:"loc,omitempty"`
	Deleted     bool      `bson:"deleted" json:"deleted"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	CreatedBy   string    `bson:"created_by" json:"created_by"`
	UpdatedBy   string    `bson:"updated_by" json:"updated_by"`
	IsAdmin     bool      `bson:"is_admin,omitempty" json:"is_admin,omitempty"`
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

func (u *User) InitWithFb(fb FBData) {
	u.FBData = fb
	u.UserName = fb.UserName
	u.FirstName = fb.FirstName
	u.LastName = fb.LastName
	u.Locale = strings.ToLower(fb.Locale)
	u.Name = fb.Name
	u.Gender = fb.Gender
}

func (u *User) InitWithGg(gg GGData) {
	u.GGData = gg
	u.FirstName = gg.Name.GivenName
	u.LastName = gg.Name.FamilyName
	u.SetName(u.FirstName, u.LastName)
	u.Gender = gg.Gender
	if gg.PlaceLived.Value != "" {
		u.City = strings.ToLower(gg.PlaceLived.Value)
	}
}

func (u *User) SetDefaults() {

	// @todo updated at should not change when use just logins
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

func (u *User) SetUserName(s string) VErrors {
	v := make(VErrors)
	n := u.Bson("UserName")
	var o User
	err := DocFindOne(bson.M{n: s}, bson.M{n: 1}, &o, 0)
	// if there is no user with such username assign username
	// else return validation error
	if err == DocNotFound {
		u.UserName = s
	} else if o.UserName != "" {
		v.Set(n, VMsg{"valid_username_taken", map[string]interface{}{"Field": n}})
	} else {
		v.Set(n, VMsg{"server_error", map[string]interface{}{}})
	}
	return v
}

func (u *User) SetBday(s string) VErrors {
	// date format layout year 2006, month 01 and day is 0
	v := make(VErrors)
	f := u.Bson("Bday")
	// layout or format of date ISO
	l := "2006-01-02"
	t, e := time.Parse(l, s)
	if e != nil {
		v.Set(f, VMsg{"valid_bday", map[string]interface{}{"Field": f}})
	} else {
		u.Bday = t
	}
	return v
}

func (u *User) Form(a, x string, t i18n.TranslateFunc) template.HTML {
	B := u.Bson
	tag := utils.Html
	type mp map[string]string
	s := template.HTML("<form action=\"" + a + "\" method=\"post\">")

	s += tag("label", mp{"for": B("FirstName"), "text": t(B("FirstName"))})
	s += tag("input", mp{"type": "text", "name": B("FirstName"), "value": u.FirstName})

	s += tag("label", mp{"for": B("LastName"), "text": t(B("LastName"))})
	s += tag("input", mp{"type": "text", "name": B("LastName"), "value": u.LastName})

	s += tag("label", mp{"for": B("Email"), "text": t(B("Email"))})
	s += tag("input", mp{"type": "email", "name": B("Email"), "value": u.Email})

	s += tag("label", mp{"for": B("Bday"), "text": t(B("Bday"))})
	s += tag("input", mp{"type": "date", "name": B("Bday"), "value": u.Bday.String()})

	s += tag("label", mp{"for": B("Gender"), "text": t(B("Gender"))})
	var cm, cf = "", ""
	if u.Gender == "male" {
		cm = "true"
	} else {
		cf = "true"
	}
	s += tag("input", mp{"type": "radio", "name": B("Gender"), "value": t("male"), "checked": cm})
	s += tag("input", mp{"type": "radio", "name": B("Gender"), "value": t("female"), "checked": cf})

	s += template.HTML(x)

	s += tag("input", mp{"type": "submit", "name": "save", "value": t("save")})

	s += template.HTML("</form>")
	return template.HTML(s)
}

func (u *User) ParseForm(m map[string][]string) {
	B := u.Bson
	// @todo think about useing reflection
	for k, v := range m {
		switch k {
		case B("UserName"):
			u.UserName = v[0]
		case B("FirstName"):
			u.FirstName = v[0]
		case B("LastName"):
			u.LastName = v[0]
		case B("Email"):
			u.Email = v[0]
		case B("Gender"):
			u.Gender = v[0]
		case B("Bday"):
			// @todo handle err properly
			_ = u.SetBday(v[0])
		}
	}
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
//@todo refactor to be dry
func (u *User) Validate(s string, bs bson.M) VErrors {
	v := Validator{}
	v.Scenario = s
	uMap := make(map[string]interface{})
	// get bson field name
	b := u.Bson
	// if validation scenario update validate
	// only provided in bson map fields
	// else validate user properties
	// do not update username
	delete(bs, b("UserName"))
	if v.Scenario == "update" {
		for k, val := range bs {
			uMap[k] = val
		}

	} else {
		uMap[b("UserName")] = u.UserName
		uMap[b("FirstName")] = u.FirstName
		uMap[b("LastName")] = u.LastName
		uMap[b("Email")] = u.Email
	}
	for k, val := range uMap {
		switch k {
		case b("UserName"):
			x := val.(string)
			v.Required(x, k)
			v.Size(x, k, 2, 100)
			v.AlphaDash(x, k)
			v.NotContainStr(x, k, []string{"admin", "administrator", "админ", "администратор"})
			v.UniqueDoc(k, u.GetC(), bson.M{k: x})

		case b("FirstName"):
			x := val.(string)
			v.Required(x, k)
			v.Size(x, k, 2, 100)

		case b("LastName"):
			x := val.(string)
			v.Required(x, k)
			v.Size(x, k, 2, 100)

		case b("Email"):
			x := val.(string)
			// it is not required validate only if not empty
			if x != "" {
				v.Email(x, k)
				v.Size(x, k, 5, 100)
			}
		}
	}

	return v.Errors
}

func (u *User) AllowBackend(id string) bool {
	var user User
	// @todo maybe cache query it will run on each adm path for each user
	err := DocFindOne(bson.M{u.Bson("ID"): bson.ObjectIdHex(id)}, bson.M{}, &user, 0)
	if err != nil {
		beego.Error("BaseCtrl.Prepare DocFindOne ", err)
		return false
	}
	r := [...]int{roleAdmin, roleEditor, roleTester}
	for _, v := range r {
		if user.Role == v {
			return true
		}
	}
	return false
}
