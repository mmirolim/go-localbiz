package ctrls

import (
	"github.com/astaxie/beego"
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"time"
)

type User struct {
	baseController
}

// @todo should let only for some roles and user himself viewable
func (c *User) Get() {
	c.TplNames = "user/user.tpl"
	// find user from id
	username := c.Ctx.Input.Param(":username")

	var user models.User
	err := models.DocFindOne(bson.M{user.Bson("UserName"): username}, bson.M{}, &user, 60)
	if err != nil {
		beego.Error(err)
		c.Abort("404")
	}
	c.Data["Uid"] = user.Id.Hex()
	c.Data["User"] = user
}

func (c *User) SignUp() {
	c.TplNames = "user/signup.tpl"
	var user models.User

	newUserData := c.GetSession("newUserData")
	if newUserData == nil {
		beego.Error("Sign up after social login for new users")
		c.Redirect("/", 302)
		return
	}

	switch newUserData.(type) {
	case models.FacebookData:
		user.InitWithFb(newUserData.(models.FacebookData))
	case models.GoogleData:
		user.InitWithGg(newUserData.(models.GoogleData))
	default:
		beego.Error("Not known social net name")
	}

	c.Data["csrfToken"] = template.HTML(c.XsrfFormHtml())
	// prefill data from social account
	c.Data["User"] = user

}
func (c *User) SignUpProc() {
	// @todo add description msg to return and aborts
	// process sign-up data on post
	if c.Ctx.Request.Method != "POST" {
		return
	}
	newUserData := c.GetSession("newUserData")
	if newUserData == nil {
		beego.Error("User.SignUpProc no newUserData")
		return
	}
	var err error
	var user models.User
	// assign Social data to user
	switch newUserData.(type) {
	case models.FacebookData:
		user.InitWithFb(newUserData.(models.FacebookData))
	case models.GoogleData:
		user.InitWithGg(newUserData.(models.GoogleData))
	default:
		beego.Error("newUserData unkown type")
		return
	}

	c.TplNames = "user/signup.tpl"

	formMap := c.Ctx.Request.PostForm
	// if user locale empty set default to current lang
	if user.Locale != "" {
		user.Locale = c.Lang
	}

	user.UserName = formMap["username"][0]
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
	if err != nil && err != models.DocNotFound {
		beego.Error("User.SignUpProc", err)
		c.Abort("500")
	}

	if existentUser.UserName != "" {
		vErrors := make(models.VErrors)
		vErrors.Set(existentUser.Bson("UserName"), models.VMsg{"valid_username_taken", map[string]interface{}{}})
		c.Data["ValidationErrors"] = vErrors
		return
	}

	vErrors, err := models.DocCreate(&user)
	if err != nil {
		beego.Error("User.SignUpProc DocCreate ", err)
		c.Abort("500")
	}
	if vErrors != nil {
		ves := make(map[string][]string)
		for k, v := range vErrors {
			for _, vmsg := range v {
				// translate field names
				vmsg.Params["Field"] = T(vmsg.Params["Field"].(string))
				msg := T(vmsg.Msg, vmsg.Params)
				ves[k] = append(ves[k], msg)
			}
		}
		c.Data["ValidationErrors"] = ves
	} else {
		// clean session
		c.DelSession("newUserData")
		//@todo should redirect after successeful signup to user account to add extra info and img
		// redirect to user's page
		err = models.DocFindOne(bson.M{"username": user.UserName}, bson.M{"username": 1}, &user, 0)

		// set user data to session
		c.SetSession("userId", user.Id.Hex())
		var urlR string
		if !check("User.SignUpProc DocFineOne ", err) {
			urlR = "/user/" + user.UserName
		} else {
			urlR = "/"
		}
		c.Redirect(urlR, 302)
		return
	}

	c.Data["csrfToken"] = template.HTML(c.XsrfFormHtml())
	c.Data["User"] = user
}

// @todo should be protected by filter only admin and user himself can edit it
func (c *User) Edit() {
	beego.Warn("User.Edit")
	var err error
	c.TplNames = "user/edit.tpl"
	// uid not nil checked in isAuth filter
	uid := c.GetSession("uid")
	var user models.User
	objId := bson.ObjectIdHex(uid.(string))
	// @todo dont cache if user edits page Or invalidate cache on update
	err = models.DocFindOne(bson.M{user.Bson("Id"): objId}, bson.M{}, &user, 0)
	if err != nil {
		beego.Error(err)
		c.Abort("404")
	}

	if c.Ctx.Request.Method == "POST" {
		formMap := c.Ctx.Request.PostForm
		flds := make(bson.M)
		B := user.Bson
		for k, v := range formMap {
			if v[0] != "" && k != "_xsrf" {
				flds[k] = v[0]
			}
		}
		// check and parse birthday
		if _, ok := flds[B("Bday")]; ok {
			flds[B("Bday")], err = time.Parse("2006-01-02", flds[B("Bday")].(string))
			if err != nil {
				delete(flds, B("Bday"))
			}
		}

		vErrors, err := models.DocUpdate(bson.M{"_id": user.Id}, &user, flds)
		//@todo handle login error properly with messages
		if err != nil {
			beego.Error("User.Edit DocUpdate ", err)
			c.Redirect("/login", 302)
			return
		}

		if vErrors != nil {
			ves := make(map[string][]string)
			for k, v := range vErrors {
				for _, vmsg := range v {
					// translate field names
					vmsg.Params["Field"] = T(vmsg.Params["Field"].(string))
					msg := T(vmsg.Msg, vmsg.Params)
					ves[k] = append(ves[k], msg)
				}
			}
			c.Data["ValidationErrors"] = ves
		} else {
			// get updated user object
			err = models.DocFindOne(bson.M{user.Bson("Id"): objId}, bson.M{}, &user, 0)
			if err != nil {
				beego.Error(err)
				c.Abort("500")
			}
		}
	}
	c.Data["csrfToken"] = template.HTML(c.XsrfFormHtml())
	c.Data["Uid"] = uid
	c.Data["User"] = user
}
