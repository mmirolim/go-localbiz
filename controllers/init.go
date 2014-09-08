package controllers

import (
	"strings"

	"github.com/beego/i18n"
	"github.com/astaxie/beego"
)

var (
	APP string
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
	this.Data["APP"] = APP

	// set default layout and layout sections
	this.Layout = "layout.tpl"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"

	// set language
	this.setLangVer()

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
func (this *baseController) setLangVer() {

	// Check URL arguments.
	lang := this.Input().Get("lang")

	// Get language info from 'Accept-Language'
	if len(lang) == 0 {
		al := this.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // only compare first 5 letters
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// Set default lang
	if len(lang) == 0 {
		lang = "ru-RU"
	}

	currentLang := langType {
		Lang : lang,
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


}

func InitApp() {
	initLocales()
}
