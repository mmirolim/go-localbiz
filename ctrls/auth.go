package ctrls

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/golang/oauth2"
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type Auth struct {
	baseController
}

var facebook, err = oauth2.NewConfig(&oauth2.Options{
	ClientID:     "892068490806056",
	ClientSecret: "194a38221d6b5b2313149b3b3510b60d",
	RedirectURL:  "http://yalp.go/auth",
},
	"https://www.facebook.com/dialog/oauth",
	"https://graph.facebook.com/oauth/access_token")

func (this *Auth) Get() {

	if err != nil {
		beego.Warn(err)
		this.Abort("403")
	}

	url := facebook.AuthCodeURL("state", "online", "auto")
	beego.Warn(url)

	this.TplNames = "login.tpl"
	this.Data["Data"] = struct {
		Url string
	}{
		url,
	}
}

func (this *Auth) Authorize() {
	this.EnableRender = false
	//@todo add msg to what was wrong
	// confirm identity
	code := this.Input().Get("code")
	beego.Warn(code)
	if code == "" {
		beego.Warn("Redirect?")
		this.Ctx.Redirect(302, "/")
		return
	}
	// exchange code to access token
	token, err := facebook.Exchange(code)
	if err != nil {
		this.Redirect("/", 302)
		return
	}
	// get public information
	t := facebook.NewTransport()
	t.SetToken(token)
	client := http.Client{Transport: t}
	resp, err := client.Get("https://graph.facebook.com/me")
	defer resp.Body.Close()
	if err != nil {
		beego.Warn(err)
		this.Redirect("/", 302)
		return
	}
	var userFbData models.FacebookData
	err = json.NewDecoder(resp.Body).Decode(&userFbData)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}
	beego.Warn(userFbData)
	var user models.User
	err = models.DocFindOne(bson.M{"fb_data.id": userFbData.Id}, bson.M{"name": 1}, &user, 60)
	if err != nil && err.Error() != "not found" {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	} else if user.UserName == "" {
		// this should be new user
		this.SetSession("newUser", 1)
		data, err := json.Marshal(&userFbData)
		if err != nil {
			beego.Error(err)
			this.Redirect("/", 302)
			return
		}
		this.SetSession("socialNet", "fb")
		this.SetSession("newUserData", data)
		this.Redirect("/signup", 302)
		return
	}

}
