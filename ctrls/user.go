package ctrls

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type User struct {
	baseController
}

func (this *User) Get() {
	beego.Warn("User.Get")
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
	if isNewUser != 1 || isNewUser == nil {
		beego.Error("Sign up after social login")
		this.Redirect("/", 302)
		return
	}
	data := this.GetSession("newUserData")
	beego.Warn(data)
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
		check("User.Get json.Unmarshal -> ", err)
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
	// process sign-up data
	if this.Ctx.Request.Method == "POST" {
		beego.Warn("this is post")
		beego.Warn(this.Ctx.Request.Method)
		formMap := this.Ctx.Request.PostForm
		beego.Warn(formMap)
		beego.Warn(formMap["first_name"][0])
	}
}
