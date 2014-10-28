package ctrls

import (
	"github.com/astaxie/beego"
	M "github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	baseController
}

// @todo should let only for some roles and user himself viewable
func (c *User) Get() {

	c.TplNames = "user/user.tpl"
	// find user from id
	n := c.Ctx.Input.Param(":username")

	var u M.User
	err := M.DocFindOne(bson.M{u.Bson("UserName"): n}, bson.M{}, &u, 60)
	if err != nil {
		beego.Error(err)
		c.Abort("404")
	}
	c.Data["uid"] = u.ID.Hex()
	c.Data["user"] = u
}

func (c *User) SignUp() {
	c.TplNames = "user/signup.tpl"
	var u M.User

	data := c.GetSession("newUserData")
	if data == nil {
		beego.Error("Sign up after social login for new users")
		c.Redirect("/", 302)
		return
	}

	switch data.(type) {
	case M.FBData:
		u.InitWithFb(data.(M.FBData))
	case M.GGData:
		u.InitWithGg(data.(M.GGData))
	default:
		beego.Error("User.SignUp Not known social net name")
	}
	c.Data["userSignUpForm"] = u.Form(UrlFor("User.SignUp"), "create", c.XsrfFormHtml(), T)

}

func (c *User) SignUpProc() {
	// @todo add description msg to return and aborts
	var u M.User
	c.TplNames = "user/signup.tpl"

	data := c.GetSession("newUserData")
	if data == nil {
		beego.Error("User.SignUpProc no newUserData")
		return
	}
	// assign Social data to user
	switch data.(type) {
	case M.FBData:
		u.InitWithFb(data.(M.FBData))
	case M.GGData:
		u.InitWithGg(data.(M.GGData))
	default:
		beego.Error("User.SignUpProc newUserData unkown type")
		return
	}

	f := c.Ctx.Request.PostForm
	// if user locale empty set default to current lang
	if u.Locale != "" {
		u.Locale = c.Lang
	}
	// parse and validate values from form
	ve := u.ParseForm(f, "create")
	c.Data["userSignUpForm"] = u.Form(UrlFor("User.SignUp"), "create", c.XsrfFormHtml(), T)
	if len(ve) != 0 {
		c.Data["vErrs"] = ve.T(T)
		return
	}
	// create new user
	ve, err := M.DocCreate(&u)
	if err != nil {
		beego.Error("User.SignUpProc DocCreate ", err)
		c.Abort("500")
	}
	if len(ve) != 0 {
		c.Data["vErrs"] = ve.T(T)
		return
	}
	// clean session
	c.DelSession("newUserData")
	//@todo should redirect after successeful signup to user account to add extra info and img

	// redirect to user's page
	n := u.Bson("UserName")
	err = M.DocFindOne(bson.M{n: u.UserName}, bson.M{n: 1}, &u, 0)

	// set user data to session
	c.SetSession("uid", u.ID.Hex())

	var r string
	if !check("User.SignUpProc DocFineOne ", err) {
		r = "/user/" + u.UserName
	} else {
		r = "/"
	}

	c.Redirect(r, 302)
}

// @todo should be protected by filter only admin and user himself can edit it
func (c *User) Edit() {
	c.TplNames = "user/edit.tpl"
	// uid not nil checked in isAuth filter
	var u M.User
	B := u.Bson
	// @todo dont cache if user edits page Or invalidate cache on update
	err := M.DocFindOne(bson.M{B("ID"): AuthUser.ID}, bson.M{}, &u, 0)
	if err != nil {
		beego.Error(err)
		c.Abort("404")
	}

	c.Data["uid"] = AuthUser.ID.Hex()
	// @todo istead hard coding use url helper
	c.Data["userEditForm"] = u.Form(UrlFor("User.Edit"), "edit", c.XsrfFormHtml(), T)
}

func (c *User) EditProc() {
	// only update own data
	c.TplNames = "user/edit.tpl"
	// uid not nil checked in isAuth filter
	var u M.User
	B := u.Bson

	f := c.Ctx.Request.PostForm
	ve := u.ParseForm(f, "edit")
	c.Data["userEditForm"] = u.Form(UrlFor("User.Edit"), "edit", c.XsrfFormHtml(), T)
	if len(ve) != 0 {
		c.Data["vErrs"] = ve.T(T)
		return
	}

	bm := M.FormToBson(f)
	// set user id who updated data
	bm[B("UpdatedBy")] = AuthUser.ID
	ve, err := M.DocUpdate(bson.M{B("ID"): AuthUser.ID}, &u, bm)
	//@todo handle login error properly with messages
	if err != nil {
		beego.Error("User.EditProc DocUpdate ", err)
		c.Abort("500")
	}
	if len(ve) != 0 {
		c.Data["vErrs"] = ve.T(T)
		return
	}

	// get updated user object
	err = M.DocFindOne(bson.M{B("ID"): AuthUser.ID}, bson.M{}, &u, 0)
	if err != nil {
		beego.Error("User.EditProc ", err)
		c.Abort("500")
	}

}
