package ctrls

import (
	"bytes"
	"encoding/gob"
	s "strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/mmirolim/beego/cache/redis"
	"github.com/nicksnyder/go-i18n/i18n"
)

var (
	APP             string
	AppVer          string
	IsPro           bool
	langTypes       map[string]string
	defaultLang     string
	cacheEnabled, _ = beego.AppConfig.Bool("cache::enabled")
	Cache, errCache = cache.NewCache("redis", `{"conn":":6379"}`)
	cachePrefix     = beego.AppConfig.String("cache::prefix") + "ctrls:"
)

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

func GetUrl(ss ...string) string {
	var u string
	// need empty first element to append first word with slash
	// ex /city/fs/one
	str := []string{""}
	for _, v := range ss {
		str = append(str, s.Replace(v, " ", "_", -1))
	}
	u = s.ToLower(s.Join(str, "/"))
	return u
}

// implement Prepare method for base router
func (this *baseController) Prepare() {

	// Setting properties
	this.Data["AppVer"] = AppVer
	this.Data["IsPro"] = IsPro
	this.Data["APP"] = APP

	// set default layout and layout sections
	this.Layout = "layout.tpl"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Header"] = "header.tpl"
	this.LayoutSections["Footer"] = "footer.tpl"

	// @todo  comment on production
	if IsPro == false {
		initLocales()
	}
	// set language
	this.setLangVer()

}

func initLocales() {
	// Init lang list
	langs := s.Split(beego.AppConfig.String("lang::types"), "|")
	names := s.Split(beego.AppConfig.String("lang::names"), "|")
	defaultLang = beego.AppConfig.String("lang::default")

	langTypes = make(map[string]string)
	for i, v := range langs {
		langTypes[v] = names[i]
	}

	i18n.MustLoadTranslationFile("./conf/en-us.all.json")
	i18n.MustLoadTranslationFile("./conf/ru-ru.all.json")

}

// set lang to use
func (this *baseController) setLangVer() {
	var lang string
	// Check URL arguments.
	urlLang := this.Input().Get("lang")

	// Get language info from 'Accept-Language'
	acceptLang := s.ToLower(this.Ctx.Request.Header.Get("Accept-Language"))

	T, err := i18n.Tfunc(urlLang, acceptLang, defaultLang)
	check("initLocales i18n.Tfunc ", err)
	// register translation func with langs
	beego.AddFuncMap("T", T)

	if langTypes[urlLang] != "" {
		lang = urlLang
	} else if langTypes[acceptLang] != "" {
		lang = acceptLang
	} else {
		lang = defaultLang
	}

	// Set lang properties
	this.Lang = lang
	this.Data["Lang"] = lang
	this.Data["CurrentLang"] = langTypes[lang]

}

func InitApp() {
	// check redis cache
	if errCache != nil {
		panic(errCache)
	}
	initLocales()
	// register getUrl func
	beego.AddFuncMap("getUrl", GetUrl)

}

func cacheIsExist(key string) bool {
	fullKey := cachePrefix + key
	return Cache.IsExist(fullKey)
}

func cachePut(key string, val interface{}, timeout int64) error {
	// serialize only structs and bytes
	// string, int, int64, float64, bool and nil will be handled by redigo
	// prepare bytes buffer
	bCache := new(bytes.Buffer)
	encCache := gob.NewEncoder(bCache)
	err := encCache.Encode(val)
	fullKey := cachePrefix + key
	Cache.Put(fullKey, bCache, timeout)
	return err
}

func cacheGet(key string, data interface{}) error {
	fullKey := cachePrefix + key
	decCache := gob.NewDecoder(bytes.NewBuffer(Cache.Get(fullKey).([]byte)))
	err := decCache.Decode(data)
	return err
}
