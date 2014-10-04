package ctrls

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"strings"
	"time"
)

type User struct {
	baseController
}

// @todo should let only for some roles and user himself viewable
func (this *User) Get() {
	this.TplNames = "user/user.tpl"
	// find user from id
	id := this.Ctx.Input.Param(":id")

	var user models.User
	objId := bson.ObjectIdHex(id)
	err := models.DocFindOne(bson.M{"_id": objId}, bson.M{}, &user, 60)
	if err != nil {
		beego.Error(err)
		this.Abort("404")
	}

	this.Data["User"] = user
}

func (this *User) SignUp() {

	this.TplNames = "user/signup.tpl"
	var user models.User

	newUserData := this.GetSession("newUserData")
	if newUserData == nil {
		beego.Error("Sign up after social login for new users")
		this.Redirect("/", 302)
		return
	}

	switch newUserData.(type) {
	case models.FacebookData:
		fb := newUserData.(models.FacebookData)
		user.UserName = fb.UserName
		user.FirstName = fb.FirstName
		user.LastName = fb.LastName
		user.Locale = strings.ToLower(fb.Locale)
		user.Name = fb.Name
		user.Gender = fb.Gender
	case models.GoogleData:
		beego.Warn("Need to implement")
	default:
		beego.Error("Not known social net name")
	}

	this.Data["csrfToken"] = template.HTML(this.XsrfFormHtml())
	// prefill data from social account
	this.Data["User"] = user

}
func (this *User) SignUpProcess() {
	// @todo add description msg to return and aborts
	// process sign-up data on post
	if this.Ctx.Request.Method != "POST" {
		return
	}
	newUserData := this.GetSession("newUserData")
	if newUserData == nil {
		return
	}
	var err error
	var user models.User
	// assign Social data to user
	switch newUserData.(type) {
	case models.FacebookData:
		user.FacebookData = newUserData.(models.FacebookData)
	case models.GoogleData:
		user.GoogleData = newUserData.(models.GoogleData)
	default:
		beego.Error("newUserData unkown type")
		return
	}
	this.TplNames = "user/signup.tpl"

	formMap := this.Ctx.Request.PostForm
	// if user locale empty set default to current lang
	if user.Locale != "" {
		user.Locale = this.Lang
	}
	user.UserName = formMap["username"][0]
	user.SetName(formMap["first_name"][0], formMap["last_name"][0])
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
	err = models.DocFindOne(bson.M{"username": user.UserName}, bson.M{"username": 1}, &existentUser, 0)
	if err != nil || err != models.DocNotFound {
		beego.Error("User.SignUpProcess", err)
		this.Abort("500")
	}
	check("User.SignUpProcess DocFindOne ", err)
	if existentUser.UserName != "" {
		vErrors := make(models.ValidationErrors)
		vErrors.Set("username", "This Username is already taken")
		this.Data["ValidationErrors"] = vErrors
		return
	}

	vErrors, err := models.DocCreate(&user)
	panicOnErr(err)

	if vErrors != nil {
		this.Data["ValidationErrors"] = vErrors
	} else {
		// clean session
		this.DelSession("newUserData")
		//@todo should redirect after successeful signup to user account to add extra info and img
		// redirect to user's page
		err = models.DocFindOne(bson.M{"username": user.UserName}, bson.M{"username": 1}, &user, 0)

		// set user data to session
		this.SetSession("userId", user.Id.Hex())
		var urlR string
		if !check("User.SignUpProcess DocFineOne ", err) {
			urlR = "/user/" + user.Id.Hex()
		} else {
			urlR = "/"
		}
		this.Redirect(urlR, 302)
		return
	}

	this.Data["csrfToken"] = template.HTML(this.XsrfFormHtml())
	this.Data["User"] = user
}
