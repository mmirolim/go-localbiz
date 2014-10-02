package ctrls

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/golang/oauth2"
	"github.com/mmirolim/yalp-go/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type Auth struct {
	baseController
}

const (
	facebook = "facebook"
	google   = "google"
)

var facebookConf, errFb = oauth2.NewConfig(&oauth2.Options{
	ClientID:     "892068490806056",
	ClientSecret: "194a38221d6b5b2313149b3b3510b60d",
	RedirectURL:  "http://yalp.go/auth",
},
	"https://www.facebook.com/dialog/oauth",
	"https://graph.facebook.com/oauth/access_token")

var googleConf, errG = oauth2.NewConfig(&oauth2.Options{
	ClientID:     "924835338434-u4osuetj2bm4j67r60g0vjur080reb97.apps.googleusercontent.com",
	ClientSecret: "DxC13zzfLWDV0TiqEMKrzkW4",
	RedirectURL:  "http://192.168.1.107.xip.io/auth",
	Scopes:       []string{"openid", "profile"},
},
	"https://accounts.google.com/o/oauth2/auth",
	"https://accounts.google.com/o/oauth2/token")

func (this *Auth) Login() {
	// check oauth2 configurations for all providers
	panicOnErr(errFb)
	panicOnErr(errG)
	isAuth := this.GetSession("isAuth")
	if isAuth == true {
		this.Redirect("/", 302)
		return
	}
	socialNet := this.Ctx.Input.Param(":socialNet")
	if socialNet == "" {
		this.TplNames = "login.tpl"
		this.Data["Data"] = struct {
			UrlFb, UrlG string
		}{
			"/login/" + facebook, "/login/" + google,
		}
		return
	}
	referer := this.Ctx.Request.Referer()
	if referer != "" {
		this.SetSession("redirectAfter", referer)
	}
	var urlR string
	switch socialNet {
	case facebook:
		urlR = facebookConf.AuthCodeURL("state", "online", "auto")
	case google:
		urlR = googleConf.AuthCodeURL("state", "online", "auto")
	default:
		urlR = "/login"
	}
	this.Redirect(urlR, 302)
}

func (this *Auth) Logout() {
	this.DestroySession()
	this.Redirect("/", 302)
}

func (this *Auth) Authorize() {
	this.EnableRender = false
	//@todo add msg to what was wrong
	// confirm identity
	code := this.Input().Get("code")
	if code == "" {
		this.Ctx.Redirect(302, "/")
		return
	}
	// exchange code to access token
	token, err := facebookConf.Exchange(code)
	if err != nil {
		this.Redirect("/", 302)
		return
	}
	// get public information
	t := facebookConf.NewTransport()
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
	var user models.User
	err = models.DocFindOne(bson.M{"fb_data.id": userFbData.Id}, bson.M{}, &user, 0)
	if err != nil && err != models.DocNotFound {
		beego.Error(err)
		this.Redirect("/login", 302)
		return
	}
	if err == models.DocNotFound {
		// this should be new user
		this.SetSession("newUser", true)
		this.SetSession("socialNet", "fb")
		this.SetSession("newUserData", userFbData)
		this.Redirect("/signup", 302)
		return
	} else {
		user.LastLoginAt = time.Now()
		vErrors, err := models.DocUpdate(bson.M{"_id": user.Id}, &user)
		//@todo handle login error properly with messages
		if err != nil {
			beego.Error(err)
			this.Redirect("/login", 302)
			return
		}
		if vErrors.Errors != nil {
			beego.Warn(vErrors)
			this.Redirect("/login", 302)
			return
		}
		redirectUrl := this.GetSession("redirectAfter")
		if redirectUrl == nil {
			redirectUrl = "/"
		}
		this.SetSession("isAuth", true)
		this.Redirect(redirectUrl.(string), 302)
	}

}
