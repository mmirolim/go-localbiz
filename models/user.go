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

// roles guest = 0, admin = 1, editor = 2, moderator = 3, tester = 4 client = 5, user = 6
type Role int

const (
	roleGuest Role = iota
	roleAdmin
	roleEditor
	roleModerator
	roleTester
	roleClient
	roleUser
)

type FBData struct {
	ID          string `bson:"id" json:"id"`
	Link        string `bson:"link" json:"link"`
	Name        string `bson:"name" json:"name"`
	FirstName   string `bson:"fname" json:"first_name"`
	LastName    string `bson:"lname" json:"last_name"`
	Gender      string `bson:"gender" json:"gender"`
	UserName    string `bson:"uname" json:"username"`
	Locale      string `bson:"locale" json:"locale"`
	AccessToken string `bson:"token"`
}

type GGData struct {
	ID         string `bson:"id" json:"id"`
	ObjectType string `bson:"objt" json:"objectType"`
	Kind       string `bson:"kind" json:"kind"`
	Etag       string `bson:"etag" json:"etag"`
	PlaceLived struct {
		Value   string `bson:"v" json:"value"`
		Primary bool   `bson:"pr" json:"primary"`
	}
	DisplayName string `bson:"dname" json:"displayName"`
	Url         string `bson:"url" json:"url"`
	Name        struct {
		FamilyName string `bson:"lname" json:"familyName"`
		GivenName  string `bson:"fname" json:"givenName"`
	} `bson:"name" json:"name"`
	Image struct {
		Url       string `bson:"url" json:"url"`
		IsDefault bool   `bson:"def" json:"isDefault"`
	} `bson:"image" json:"image"`
	FirstName   string `bson:"fname" json:"first_name"`
	LastName    string `bson:"lame" json:"last_name"`
	Gender      string `bson:"gender" json:"gender"`
	UserName    string `bson:"uname" json:"username"`
	Lang        string `bson:"lang" json:"language"`
	isPlusUser  string `bosn:"pluser" json:"isPlusUser"`
	Verified    bool   `bson:"vrd" json:"verified"`
	AccessToken string `bson:"token"`
}

// @important  if Field naming changes validation also should be changed
//@todo fix validation for email now requiring
type User struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	UserName    string        `bson:"uname" json:"username"`
	Email       string        `bson:"email" json:"email"`
	Name        string        `bson:"name" json:"name"`
	City        string        `bson:"city" json:"city"`
	FirstName   string        `bson:"fname" json:"first_name"`
	LastName    string        `bson:"lname" json:"last_name"`
	Gender      string        `bson:"gender" json:"gender"`
	Locale      string        `bson:"locale" json:"locale" json:"locale"`
	LastLoginAt time.Time     `bson:"ll_at" json:"last_login_at"`
	Role        Role          `bson:"role" json:"role"`
	Bday        time.Time     `bson:"bday,omitempty" json:"bday,omitempty"`
	FBData      `bson:"fbd" json:"fb_data"`
	GGData      `bson:"ggd" json:"gg_data"`
	Address     `bson:"addr" json:"address"`
	Geo         `bson:"loc,omitempty" json:"loc,omitempty"`
	Deleted     bool          `bson:"del" json:"deleted"`
	UpdatedAt   time.Time     `bson:"up_at" json:"updated_at"`
	CreatedAt   time.Time     `bson:"cr_at" json:"created_at"`
	CreatedBy   bson.ObjectId `bson:"cr_by,omitempty" json:"created_by"`
	UpdatedBy   bson.ObjectId `bson:"up_by,omitempty" json:"updated_by"`
	IsAdmin     bool          `bson:"adm,omitempty" json:"is_admin,omitempty"`
}

func (u User) GetC() string {
	return "users"
}

func (u User) GetIndex() []mgo.Index {
	// define indexes
	ub := u.Bson
	ab := Address{}.Bson
	return []mgo.Index{
		mgo.Index{
			Key:    []string{ub("UserName")},
			Unique: true,
		},
		mgo.Index{
			Key: []string{ub("Name")},
		},
		mgo.Index{
			Key:    []string{ub("FBData") + ".id"},
			Unique: true,
		},
		mgo.Index{
			Key:    []string{ub("GGData") + ".id"},
			Unique: true,
		},
		mgo.Index{
			Key: []string{ub("Gender")},
		},

		mgo.Index{
			Key: []string{ub("Role")},
		},
		mgo.Index{
			Key: []string{ub("Bday")},
		},

		mgo.Index{
			Key:    []string{ub("Email")},
			Unique: true,
		},

		mgo.Index{
			Key: []string{ub("LastLoginAt")},
		},

		mgo.Index{
			Key: []string{"$2dsphere:" + ub("Geo")},
		},
		mgo.Index{
			Key: []string{ub("Address") + "." + ab("City"), ub("Name")},
		},
		mgo.Index{
			Key: []string{ub("Address") + "." + ab("District"), ub("Name")},
		},
		mgo.Index{
			Key: []string{ub("Deleted")},
		},
		mgo.Index{
			Key: []string{ub("UpdatedBy")},
		},
		mgo.Index{
			Key: []string{ub("UpdatedAt")},
		},
		mgo.Index{
			Key: []string{ub("CreatedAt")},
		},
		mgo.Index{
			Key: []string{ub("CreatedBy")},
		},
	}
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

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.LastLoginAt.IsZero() {
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

// a - action, s - scenario like "create, edit", x - xsrf token, t - i18n translation func
func (u *User) Form(a, s, x string, t i18n.TranslateFunc) template.HTML {
	B := u.Bson
	tag := utils.Html
	type mp map[string]string
	f := template.HTML("<form action=\"" + a + "\" method=\"post\">")
	// add username field if user created
	if s == "create" {
		f += tag("label", mp{"for": B("UserName"), "text": t(B("UserName"))})
		f += tag("input", mp{"type": "text", "name": B("UserName"), "value": u.UserName})
	} else {
		// add hidden user id field
		f += tag("input", mp{"type": "hidden", "name": B("ID"), "value": u.ID.Hex()})
	}
	f += tag("label", mp{"for": B("FirstName"), "text": t(B("FirstName"))})
	f += tag("input", mp{"type": "text", "name": B("FirstName"), "value": u.FirstName})

	f += tag("label", mp{"for": B("LastName"), "text": t(B("LastName"))})
	f += tag("input", mp{"type": "text", "name": B("LastName"), "value": u.LastName})

	f += tag("label", mp{"for": B("Email"), "text": t(B("Email"))})
	f += tag("input", mp{"type": "email", "name": B("Email"), "value": u.Email})

	f += tag("label", mp{"for": B("Bday"), "text": t(B("Bday"))})
	f += tag("input", mp{"type": "date", "name": B("Bday"), "value": u.Bday.String()})

	f += tag("label", mp{"for": B("Gender"), "text": t(B("Gender"))})
	var cm, cf = "", ""
	if u.Gender == "male" {
		cm = "true"
	} else {
		cf = "true"
	}
	f += template.HTML("<select name=\"" + B("Gender") + "\">")
	f += tag("option", mp{"value": t("male"), "selected": cm, "text": t("male")})
	f += tag("option", mp{"value": t("female"), "selected": cf, "text": t("female")})
	f += template.HTML("</select>")
	f += template.HTML(x)

	f += tag("input", mp{"type": "submit", "name": "save", "value": t("save")})

	f += template.HTML("</form>")
	return template.HTML(f)
}

// parse form to user struct and validate with s - scenario
func (u *User) ParseForm(m map[string][]string, s string) VErrors {
	B := u.Bson
	vErrs := make(VErrors)
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
			// if bday field not empty convert to proper type
			if v[0] != "" {
				e := u.SetBday(v[0])
				if e != nil {
					vErrs.Add(e)
				}
			}
		}
	}
	// create bson field map to validate
	b := FormToBson(m)
	e := u.Validate(s, b)
	if e != nil {
		vErrs.Add(e)
	}
	return vErrs
}

// validate User properties with s - scenario
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
	if v.Scenario == "create" {
		uMap[b("UserName")] = u.UserName
		uMap[b("FirstName")] = u.FirstName
		uMap[b("LastName")] = u.LastName
		uMap[b("Email")] = u.Email
	} else {
		delete(bs, b("UserName"))
		for k, val := range bs {
			uMap[k] = val
		}
	}
	for k, val := range uMap {
		switch k {
		case b("UserName"):
			x := val.(string)
			v.Required(x, k)
			v.Size(x, k, 2, 100)
			v.AlphaDash(x, k)
			v.NotContainStr(x, k, []string{"admin", "administrator", "админ", "администратор"})
			v.UniqueDoc(k, x, u.GetC())

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
			// it is not required, validate only if not empty
			if x == "" {
				break
			}
			v.Email(x, k)
			v.Size(x, k, 5, 100)
			v.UniqueDoc(k, x, u.GetC())
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
	r := [...]Role{roleAdmin, roleEditor, roleTester}
	for _, v := range r {
		if user.Role == v {
			return true
		}
	}
	return false
}

// get bson tag from field name
func (u User) Bson(f string) string {
	return dic.User.Bson(f)
}

// get field name from bson tag
func (u User) Field(f string) string {
	return dic.User.Field(f)
}

func (u *User) Create() (VErrors, error) {
	// set unique ObjectID
	if u.ID.Hex() == "" {
		u.ID = bson.NewObjectId()
	}
	return DocCreate(u)
}

func (u *User) Update(f bson.M) (VErrors, error) {
	return DocUpdate(bson.M{"_id": u.ID}, u, f)
}

func (u *User) FindById(id string, t int64) error {
	return DocFindOne(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{}, u, t)
}

func (u *User) FindOne(q bson.M, t int64) error {
	return DocFindOne(q, bson.M{}, u, t)
}
