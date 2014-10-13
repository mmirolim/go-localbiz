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
	uid := c.GetSession("uid")
	// if user authenticated redirect
	if uid != nil {
		c.Redirect("/", 302)
		return
	}
	socialNet := c.Ctx.Input.Param(":socialNet")
	if socialNet == "" {
		c.TplNames = "login.tpl"
		loginUrl := "/login/"
		c.Data["Data"] = struct {
			UrlFb, UrlG string
		}{
			loginUrl + facebook, loginUrl + google,
		}
		return
	}
	// store referrer to redirect to a page where user logged in
	referer := c.Ctx.Request.Referer()
	if referer != "" {
		c.SetSession("redirectAfter", referer)
	}
	// @todo add csrf tokens as state
	var urlR string
	state := c.XsrfToken()
	switch socialNet {
	case facebook:
		urlR = facebookConf.AuthCodeURL(state, "online", "auto")
	case google:
		urlR = googleConf.AuthCodeURL(state, "online", "auto")
	default:
		urlR = "/login"
	}
	c.Redirect(urlR, 302)
}

func (c *Auth) Logout() {
	c.DestroySession()
	c.Redirect("/", 302)
}

func (c *Auth) Authorize() {
	c.EnableRender = false
	//@todo add msg to what was wrong
	socialNet := c.Ctx.Input.Param(":socialNet")
	// confirm identity
	state := c.Input().Get("state")
	if state != c.XsrfToken() {
		c.Abort("403")
	}
	code := c.Input().Get("code")
	if code == "" || socialNet == "" {
		c.Ctx.Redirect(302, "/")
		return
	}
	// declare var required for oauth providers
	var providerConf *oauth2.Config
	var userInfoUrl string
	var userData interface{}
	var userFbData models.FacebookData
	var userGgData models.GoogleData
	switch socialNet {
	case facebook:
		userData = &userFbData
		providerConf = facebookConf
		userInfoUrl = "https://graph.facebook.com/me"
	case google:
		userData = &userGgData
		providerConf = googleConf
		userInfoUrl = "https://www.googleapis.com/plus/v1/people/me"
	default:
		c.Abort("400")
	}
	// exchange code to access token
	token, err := providerConf.Exchange(code)
	if err != nil {
		c.Redirect("/", 302)
		return
	}
	// get public information
	t := providerConf.NewTransport()
	t.SetToken(token)
	client := http.Client{Transport: t}
	resp, err := client.Get(userInfoUrl)
	defer resp.Body.Close()

	if err != nil {
		beego.Warn(err)
		c.Redirect("/", 302)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(userData)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}

	var user models.User
	var socialId bson.M
	// search by social id
	switch userData.(type) {
	case *models.FacebookData:
		socialId = bson.M{"fb_data.id": userData.(*models.FacebookData).Id}
		// get value and typecast to proper data type
		c.SetSession("newUserData", *userData.(*models.FacebookData))
	case *models.GoogleData:
		socialId = bson.M{"gg_data.id": userData.(*models.GoogleData).Id}
		c.SetSession("newUserData", *userData.(*models.GoogleData))
	default:
		beego.Error("userData type unkown")
	}
	err = models.DocFindOne(socialId, bson.M{}, &user, 0)
	if err != nil && err != models.DocNotFound {
		beego.Error(err)
		c.Redirect("/login", 302)
		return
	}
	if err == models.DocNotFound {
		// c should be new user
		c.Redirect("/signup", 302)
		return
	} else {
		// delete newUserData if existent user
		c.DelSession("newUserData")
		// update last login
		q := bson.M{user.Bson("Id"): user.Id}
		fld := bson.M{user.Bson("LastLoginAt"): time.Now()}
		vErrors, err := models.DocUpdate(q, &user, fld)
		//@todo handle login error properly with messages
		if err != nil {
			beego.Error("Auth.Authorize DocUpdate ", err)
			c.Redirect("/login", 302)
			return
		}
		if vErrors != nil {
			beego.Warn(vErrors)
			c.Redirect("/login", 302)
			return
		}
		redirectUrl := c.GetSession("redirectAfter")
		if redirectUrl == nil {
			redirectUrl = "/"
		}
		c.SetSession("uid", user.Id.Hex())
		c.Redirect(redirectUrl.(string), 302)
	}

}
