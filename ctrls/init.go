package ctrls

import (
	"strings"

	"github.com/astaxie/beego"
	M "github.com/mmirolim/yalp-go/models"
	"github.com/nicksnyder/go-i18n/i18n"
	"gopkg.in/mgo.v2/bson"
)

var (
	APP       string
	AppVer    string
	IsPro     bool
	langTypes map[string]string
	dLang = beego.AppConfig.String("lang::default") // set default lang
	T         i18n.TranslateFunc
	AuthUser  M.User
	UrlFor = beego.UrlFor
)

func InitApp() {
	initLocales()
	// register getUrl func
	// @todo use default url builder
	beego.AddFuncMap("getUrl", GetUrl)
}

// implement Prepare method for base router
func (c *baseController) Prepare() {
	// Setting properties
	c.Data["AppVer"] = AppVer
	c.Data["IsPro"] = IsPro
	c.Data["APP"] = APP

	// set default layout and layout sections
	c.Layout = "layout.tpl"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["Header"] = "header.tpl"
	c.LayoutSections["Footer"] = "footer.tpl"

	// @todo  comment on production
	if IsPro == false {
		initLocales()
	}
	// set language
	c.setLangVer()

	// set Authenticated User available for all controllers
	uid := c.GetSession("uid")
	if uid != nil {
		err := M.DocFindOne(bson.M{AuthUser.Bson("ID"): bson.ObjectIdHex(uid.(string))}, bson.M{}, &AuthUser, 0)
		if err != nil {
			beego.Error("BaseCtrl.Prepare DocFindOne ", err)
			c.Abort("500")
		}
	}
	c.Data["AuthUser"] = AuthUser
}

func initLocales() {
	// Init lang list
	langs := strings.Split(beego.AppConfig.String("lang::types"), "|")
	names := strings.Split(beego.AppConfig.String("lang::names"), "|")

	langTypes = make(map[string]string)
	for i, v := range langs {
		langTypes[v] = names[i]
	}

	i18n.MustLoadTranslationFile("./conf/en-us.all.json")
	i18n.MustLoadTranslationFile("./conf/ru-ru.all.json")

}

// set lang to use
func (c *baseController) setLangVer() {
	var lang string
	// Check URL arguments if no set as default.
	uLang := strings.ToLower(c.Input().Get("lang"))
	if uLang == "" {
		uLang = dLang
	}
	// Get language info from 'Accept-Language'
	aLang := strings.ToLower(c.Ctx.Request.Header.Get("Accept-Language"))
	if len(aLang) > 4 && langTypes[aLang[:5]] != "" {
		aLang = aLang[:5]
	} else {
		aLang = ""
	}
	Tfn, err := i18n.Tfunc(uLang, aLang, dLang)
	// set T func in ctrls
	T = Tfn
	check("initLocales i18n.Tfunc ", err)
	// register translation func with langs
	beego.AddFuncMap("T", Tfn)

	switch {
	case langTypes[uLang] != "":
		lang = uLang
	case langTypes[aLang] != "":
		lang = aLang
	default:
		lang = dLang
	}

	// Set lang properties
	c.Lang = lang
	c.Data["Lang"] = lang
	c.Data["CurrentLang"] = langTypes[lang]
}

// convenience type for i18n
type tm map[string]interface{}

// base router with global settings for all other routers
type baseController struct {
	beego.Controller
	Lang string
}
// @todo should not be used instead use beego UrlFor
func GetUrl(ss ...string) string {
	var u string
	// need empty first element to append first word with slash
	// ex /city/fs/one
	str := []string{""}
	for _, v := range ss {
		str = append(str, strings.Replace(v, " ", "_", -1))
	}
	u = strings.ToLower(strings.Join(str, "/"))
	return u
}

func check(s string, e error) bool {
	if e != nil {
		beego.Error(s + e.Error())
		return true
	}

	return false
}

func panicOnErr(e error) {
	if e != nil {
		panic(e)
	}
}
