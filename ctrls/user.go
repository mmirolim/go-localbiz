package ctrls

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type User struct {
	baseController
}

func (this *User) Get() {
	this.TplNames = "user/user.tpl"
	// find user from id
	id := this.Ctx.Input.Param(":id")
	var user models.User
	err := models.DocFindOne(bson.M{"_id": id}, bson.M{}, &user, 60)
	if err != nil {
		beego.Error(err)
		this.Abort("404")
	}
	isNewUser := this.GetSession("isNewUser").(bool)
	isAuthenticated := this.GetSession("authenticated").(bool)
	userFromSess := this.GetSession("user")
	this.Data["User"] = user
	this.Data["IsNew"] = isNewUser
	beego.Warn(userFromSess)
	beego.Warn(isAuthenticated)
}

func (this *User) SignUp() {
	this.TplNames = "user/signup.tpl"
	var user models.User
	isNewUser := this.GetSession("newUser")
	if isNewUser != true || isNewUser == nil {
		beego.Error("Sign up after social login")
		this.Redirect("/", 302)
		return
	}
	data := this.GetSession("newUserData")
	if data == nil {
		beego.Error("newUserData should not be nil")
		this.Redirect("/", 302)
		return
	}
	socialNet := this.GetSession("socialNet")
	if socialNet == nil {
		beego.Error("socialNet should not be nil")
		this.Redirect("/", 302)
		return
	}
	var fbData models.FacebookData

	switch socialNet {
	case "fb":
		err := json.Unmarshal(data.([]byte), &fbData)
		// @todo maybe should inform and redirect
		check("User.Get json.Unmarshal -> ", err)
		user.FacebookData = fbData
		user.UserName = fbData.UserName
		user.FirstName = fbData.FirstName
		user.LastName = fbData.LastName
		user.Locale = strings.ToLower(fbData.Locale)
		user.Name = fbData.Name
		user.Gender = fbData.Gender
	default:
		beego.Error("Not known social net name")
	}

	// prefill data from social account
	this.Data["User"] = user

	// process sign-up data from form
	if this.Ctx.Request.Method != "POST" {
		return
	}

	formMap := this.Ctx.Request.PostForm
	user.UserName = formMap["username"][0]
	user.Name = formMap["first_name"][0] + formMap["last_name"][0]
	user.FirstName = formMap["first_name"][0]
	user.LastName = formMap["last_name"][0]
	user.Email = formMap["email"][0]
	// date format layout year 2006, month 01 and day is 02
	bday := formMap["bday"][0]
	if bday != "" {
		user.Bday, err = time.Parse("2006-01-02", bday)
		check("User.SignUp Bday format error ->", err)
	}
	user.Gender = formMap["gender"][0]
	// check if username is free
	var existentUser models.User
	err := models.DocFindOne(bson.M{"username": user.UserName}, bson.M{"username": 1}, &existentUser, 0)
	if existentUser.UserName != "" {
		this.Data["ValidationErrors"] = []struct {
			Key     string
			Message string
		}{
			{"Username", "This username is already taken"},
		}
		return
	}
	vErrors, err := models.DocCreate(&user)
	panicOnErr(err)
	if vErrors != nil {
		this.Data["ValidationErrors"] = vErrors
		beego.Warn(vErrors)
	}
	// clean session
	this.DelSession("newUser")
	this.DelSession("socialNet")
	this.DelSession("newUserData")

}
