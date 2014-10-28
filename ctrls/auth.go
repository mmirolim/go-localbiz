package ctrls

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/golang/oauth2"
	M "github.com/mmirolim/yalp-go/models"
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
	RedirectURL:  "http://yalp.go/auth/" + facebook,
},
	"https://www.facebook.com/dialog/oauth",
	"https://graph.facebook.com/oauth/access_token")

var googleConf, errG = oauth2.NewConfig(&oauth2.Options{
	ClientID:     "924835338434-u4osuetj2bm4j67r60g0vjur080reb97.apps.googleusercontent.com",
	ClientSecret: "DxC13zzfLWDV0TiqEMKrzkW4",
	RedirectURL:  "http://192.168.1.107.xip.io/auth/" + google,
	Scopes:       []string{"openid", "profile"},
},
	"https://accounts.google.com/o/oauth2/auth",
	"https://accounts.google.com/o/oauth2/token")

func (c *Auth) Login() {
	// check oauth2 configurations for all providers
	panicOnErr(errFb)
	panicOnErr(errG)
	// if user authenticated redirect
	if AuthUser.ID.Hex() != "" {
		c.Redirect("/", 302)
		return
	}
	prv := c.Ctx.Input.Param(":provider")
	if prv == "" {
		c.TplNames = "login.tpl"
		l := "/login/"
		c.Data["url"] = struct {
			FB, GG string
		}{
			l + facebook, l + google,
		}
		return
	}
	// store referrer to redirect to a page where user logged in
	referer := c.Ctx.Request.Referer()
	if referer != "" {
		c.SetSession("redirectAfter", referer)
	}
	// @todo add csrf tokens as state
	var url string
	s := c.XsrfToken()
	switch prv {
	case facebook:
		url = facebookConf.AuthCodeURL(s, "online", "auto")
	case google:
		url = googleConf.AuthCodeURL(s, "online", "auto")
	default:
		url = "/login"
	}

	c.Redirect(url, 302)
}

func (c *Auth) Logout() {
	// empty data in Authenticated user struct
	AuthUser = M.User{}
	c.DestroySession()
	c.Redirect("/", 302)
}

func (c *Auth) Authorize() {
	c.EnableRender = false
	//@todo add msg to what was wrong
	prv := c.Ctx.Input.Param(":provider")
	// confirm identity
	state := c.Input().Get("state")
	// @todo remove social net check
	// temp fix for google login
	if state != c.XsrfToken() && prv != "google" {
		beego.Warn("Auth.Authorize state mismatch")
		c.Abort("403")
	}
	code := c.Input().Get("code")
	if code == "" || prv == "" {
		c.Ctx.Redirect(302, "/")
		return
	}
	// declare var required for oauth providers
	var pConf *oauth2.Config
	var userURL string
	var uData interface{}
	var uFBData M.FBData
	var uGGData M.GGData
	switch prv {
	case facebook:
		uData = &uFBData
		pConf = facebookConf
		userURL = "https://graph.facebook.com/me"
	case google:
		uData = &uGGData
		pConf = googleConf
		userURL = "https://www.googleapis.com/plus/v1/people/me"
	default:
		c.Abort("400")
	}
	// exchange code to access token
	token, err := pConf.Exchange(code)
	if err != nil {
		c.Redirect("/", 302)
		return
	}
	// get public information
	t := pConf.NewTransport()
	t.SetToken(token)
	client := http.Client{Transport: t}
	r, err := client.Get(userURL)
	defer r.Body.Close()

	if err != nil {
		beego.Warn(err)
		c.Redirect("/", 302)
		return
	}

	err = json.NewDecoder(r.Body).Decode(uData)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}

	var user M.User
	var sid bson.M

	// search by social id
	switch uData.(type) {
	case *M.FBData:
		sid = bson.M{user.Bson("FBData") + ".id": uData.(*M.FBData).ID}
		// get value and typecast to proper data type
		c.SetSession("newUserData", *uData.(*M.FBData))
	case *M.GGData:
		sid = bson.M{user.Bson("GGData") + ".id": uData.(*M.GGData).ID}
		c.SetSession("newUserData", *uData.(*M.GGData))
	default:
		beego.Error("userData type unkown")
	}

	// find use by social Id used to login
	err = M.DocFindOne(sid, bson.M{}, &user, 0)
	switch {
	case err != nil && err != M.DocNotFound:
		beego.Error(err)
		c.Redirect("/login", 302)
		return
	case err == M.DocNotFound:
		// c should be new user
		c.Redirect("/signup", 302)
		return
	}

	// delete newUserData if existent user
	c.DelSession("newUserData")
	// update last login
	q := bson.M{user.Bson("ID"): user.ID}
	f := bson.M{user.Bson("LastLoginAt"): time.Now()}
	// udpate user last login time
	verrs, err := M.DocUpdate(q, &user, f)
	//@todo handle login error properly with messages
	if err != nil {
		beego.Error("Auth.Authorize DocUpdate ", err)
		c.Redirect("/login", 302)
		return
	}
	if verrs != nil {
		beego.Warn(verrs)
		c.Redirect("/login", 302)
		return
	}
	rURL := c.GetSession("redirectAfter")
	if rURL == nil {
		rURL = "/"
	}
	c.SetSession("uid", user.ID.Hex())
	c.Redirect(rURL.(string), 302)

}
