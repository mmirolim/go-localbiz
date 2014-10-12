package ctrls

import (
	"strings"

	"github.com/astaxie/beego"
	"github.com/nicksnyder/go-i18n/i18n"
)

var (
	APP         string
	AppVer      string
	IsPro       bool
	langTypes   map[string]string
	defaultLang string
	T           i18n.TranslateFunc
)

// convenience type for i18n
type tm map[string]interface{}

// base router with global settings for all other routers
type baseController struct {
	beego.Controller
	Lang string
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

}

func initLocales() {
	// Init lang list
	langs := strings.Split(beego.AppConfig.String("lang::types"), "|")
	names := strings.Split(beego.AppConfig.String("lang::names"), "|")
	defaultLang = beego.AppConfig.String("lang::default")

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
	// Check URL arguments.
	urlLang := strings.ToLower(c.Input().Get("lang"))

	// Get language info from 'Accept-Language'
	acceptLang := strings.ToLower(c.Ctx.Request.Header.Get("Accept-Language"))

	Tfn, err := i18n.Tfunc(urlLang, acceptLang, defaultLang)
	// set T func in ctrls
	T = Tfn
	check("initLocales i18n.Tfunc ", err)
	// register translation func with langs
	beego.AddFuncMap("T", Tfn)
	if langTypes[urlLang] != "" {
		lang = urlLang
	} else if langTypes[acceptLang] != "" {
		lang = acceptLang
	} else {
		lang = defaultLang
	}

	// Set lang properties
	c.Lang = lang
	c.Data["Lang"] = lang
	c.Data["CurrentLang"] = langTypes[lang]
}

func InitApp() {
	initLocales()
	// register getUrl func
	beego.AddFuncMap("getUrl", GetUrl)
}
