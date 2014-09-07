package controllers

import (
	"strings"

	"github.com/beego/i18n"
	"github.com/astaxie/beego"
)

var (
	AppVer string
	IsPro bool
	langTypes []*langType
)
// languages struct
type langType struct {
	Lang, Name string
}
// base router with global settings for all other routers
type baseController struct {
	beego.Controller
	i18n.Locale
}

// implement Prepare method for base router
func (this *baseController) Prepare() {
	// Setting properties
	this.Data["AppVer"] = AppVer
	this.Data["IsPro"]  = IsPro
	// Redirect to make URL clean
	if this.setLangVer() {
		i := strings.Index(this.Ctx.Request.RequestURI, "?")
		this.Redirect(this.Ctx.Request.RequestURI[:i], 302)
		return
	}

}

func initLocales() {
	// Init lang list
	langs := strings.Split(beego.AppConfig.String("lang::types"), "|")
	names := strings.Split(beego.AppConfig.String("lang::names"), "|")
	langTypes = make([]*langType, 0, len(langs))

	for i, v := range langs {
		langTypes = append(langTypes, &langType{ Lang: v, Name: names[i]})
	}

	for _, lang := range langs {
		beego.Trace("Loading language" + lang)
		if err := i18n.SetMessage(lang, "conf/locale_" + lang + ".ini"); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			return
		}
	}

}

// set lang to use
func (this *baseController) setLangVer() bool {
	var isNeedRedir, hasCookie = false, false

	// 1. Check URL arguments.
	lang := this.Input().Get("lang")

	// 2. Get lang information from cookies
	if len(lang) == 0 {
		lang = this.Ctx.GetCookie("lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check if lang in cookie exists
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir, hasCookie = false, false
	}

	// 3. Get language info from 'Accept-Language'
	if len(lang) == 0 {
		al := this.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // only compare first 5 letters
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. Set default lang
	if len(lang) == 0 {
		lang = "ru-RU"
		isNeedRedir = false
	}

	currentLang := langType {
		Lang : lang,
	}

	// Save lang in cookies
	if !hasCookie {
		this.Ctx.SetCookie("lang", currentLang.Lang, 1 << 31-1, "/")
	}

	restLangs := make([]*langType, 0, len(langTypes) - 1)
	for _, v := range langTypes {
		if lang != v.Lang {
			restLangs = append(restLangs, v)
		} else {
			currentLang.Name = v.Name
		}
	}

	// Set lang properties
	this.Lang = lang
	this.Data["Lang"] = currentLang.Lang
	this.Data["CurrentLang"] = currentLang.Name
	this.Data["ResLangs"] = restLangs

	// return if redirect needed
	return isNeedRedir
}

func InitApp() {
	initLocales()
}
