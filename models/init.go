package models

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/nicksnyder/go-i18n/i18n"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"regexp"
	"strings"
)

// connection
var (
	MongoHost       string
	MongoDbName     string
	MgoSession      *mgo.Session
	CacheEnabled, _ = beego.AppConfig.Bool("cache::enabled")
	Cache, errCache = cache.NewCache("redis", `{"conn":":6379"}`)
	cachePrefix     = beego.AppConfig.String("cache::prefix")
	DocNotFound     = mgo.ErrNotFound
	FieldDic        map[string]map[string]map[string]string
	T               = i18n.IdentityTfunc
	// regex patters
	reg_alpha_dash = regexp.MustCompile("[^\\d\\w-_]")
	reg_email      = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")
)

type DocModel interface {
	GetC() string
	GetIndex() []mgo.Index
	FmtFields()
	SetDefaults()
	Validate(s string, bs bson.M) VErrors
	GetLocation() Geo
}

func InitConnection() {
	MongoHost = beego.AppConfig.String("db::mongohost")
	MongoDbName = beego.AppConfig.String("db::mongodbname")
	session, err := mgo.Dial(MongoHost)
	if err != nil {
		beego.Error("Could not connect to mongo instance", err)
	}
	// import to pkg scope
	MgoSession = session

	// all models
	var user User
	var fs FoodService
	var addr Address
	FieldDic = make(map[string]map[string]map[string]string)
	FieldDic["User"] = make(map[string]map[string]string)
	FieldDic["Address"] = make(map[string]map[string]string)
	FieldDic["FoodService"] = make(map[string]map[string]string)
	// define maps bsonToUser and userToBson field names
	FieldDic["User"]["FieldBson"] = FieldBsonDic(&user)
	FieldDic["User"]["BsonField"] = BsonFieldDic(&user)
	FieldDic["Address"]["FieldBson"] = FieldBsonDic(&addr)
	FieldDic["FoodService"]["FieldBson"] = FieldBsonDic(&fs)

	// init indexes of models and panic if something wrong
	err = DocInitIndex(&fs)
	panicOnErr(err)
	err = DocInitIndex(&user)
	panicOnErr(err)
	// check new cache
	panicOnErr(errCache)

	// register Model structs for gob encoding
	gob.Register(user)
	gob.Register(fs)
	gob.Register(addr)
	gob.Register(FBData{})
	gob.Register(GGData{})

}

func DocInitIndex(m DocModel) error {
	var err error
	sess := MgoSession.Copy()
	defer sess.Close()
	for _, v := range m.GetIndex() {
		err = sess.DB(MongoDbName).C(m.GetC()).EnsureIndex(v)
		if check("InitIndex -> ", err) {
			break
		}
	}

	return err
}

func DocFind(q bson.D, f bson.M, m DocModel, data interface{}, timeout int64) error {

	var err error
	cacheKey := genCacheKey(m.GetC(), "DocFind", q, f)
	if CacheEnabled && cacheIsExist(cacheKey) && timeout > 0 {
		cacheGet(cacheKey, data)
		return err
	}
	sess := MgoSession.Copy()
	defer sess.Close()

	collection := sess.DB(MongoDbName).C(m.GetC())
	// limit is important when all used, may consume all memory
	// @todo maybe memory consumption reduces if not all fields retrieved?
	iter := collection.Find(q).Select(f).Limit(5000).Iter()
	err = iter.All(data)
	if CacheEnabled && err == nil && timeout > 0 {
		cachePut(cacheKey, data, timeout)
	}

	return err
}

func DocFindOne(q bson.M, f bson.M, m DocModel, timeout int64) error {

	var err error
	cacheKey := genCacheKey(m.GetC(), "DocFindOne", q, f)
	if CacheEnabled && cacheIsExist(cacheKey) && timeout > 0 {
		cacheGet(cacheKey, m)
		return err
	}
	sess := MgoSession.Copy()
	defer sess.Close()

	collection := sess.DB(MongoDbName).C(m.GetC())
	err = collection.Find(q).Select(f).One(m)

	if CacheEnabled && err == nil && timeout > 0 {
		cachePut(cacheKey, m, timeout)
	}

	return err
}

// @todo refactor maybe loc.coor should be passed by f?
// currently mongo 2.6.4 does not support geoNear with subset of fields
func DocFindNear(min, max int, m DocModel, data interface{}, timeout int64) error {

	var err error
	q := bson.D{
		{"geoNear", m.GetC()},
		{"near", bson.D{{"type", "Point"}, {"coordinates", m.GetLocation().Coordinates}}},
		{"spherical", true},
		{"minDistance", min},
		{"maxDistance", max},
	}
	cacheKey := genCacheKey(m.GetC(), "DocFindNear", q)
	if CacheEnabled && cacheIsExist(cacheKey) && timeout > 0 {
		cacheGet(cacheKey, data)
		return err
	}

	sess := MgoSession.Copy()
	defer sess.Close()

	err = sess.DB(MongoDbName).Run(q, data)

	if CacheEnabled && err == nil && timeout > 0 {
		cachePut(cacheKey, data, timeout)
	}

	return err
}

// find all distinct tags in arrays and count docs with each tag
func DocCountDistinct(m DocModel, match bson.M, category string, data interface{}, timeout int64) error {

	var err error
	q := []bson.M{
		{"$match": match},
		{"$project": bson.M{category: 1}},
		{"$unwind": "$" + category},
		{"$group": bson.D{{"_id", "$" + category}, {"count", bson.M{"$sum": 1}}}},
	}
	cacheKey := genCacheKey(m.GetC(), "DocCountDistinct", q)
	if CacheEnabled && cacheIsExist(cacheKey) && timeout > 0 {
		cacheGet(cacheKey, data)
		return err
	}

	sess := MgoSession.Copy()
	defer sess.Close()

	collection := sess.DB(MongoDbName).C(m.GetC())
	pipe := collection.Pipe(q)
	err = pipe.Iter().All(data)

	if CacheEnabled && err == nil && timeout > 0 {
		cachePut(cacheKey, data, timeout)
	}

	return err
}

func DocCreate(m DocModel) (VErrors, error) {
	var err error
	// validate model before inserting
	vErrors := m.Validate("create", bson.M{})
	if vErrors != nil {
		return vErrors, err
	}

	sess := MgoSession.Copy()
	defer sess.Close()

	// set defaults
	m.FmtFields()
	m.SetDefaults()

	collection := sess.DB(MongoDbName).C(m.GetC())
	err = collection.Insert(m)

	return vErrors, err
}

func DocUpdate(q bson.M, m DocModel, flds bson.M) (VErrors, error) {
	var err error
	// validate model before inserting
	vErrors := m.Validate("update", flds)
	if vErrors != nil {
		return vErrors, err
	}

	sess := MgoSession.Copy()
	defer sess.Close()

	// set defaults
	m.FmtFields()
	m.SetDefaults()

	collection := sess.DB(MongoDbName).C(m.GetC())
	// update fields with $set
	err = collection.Update(q, bson.M{"$set": flds})

	return vErrors, err

}

type Validator struct {
	Scenario string // should be create, update
	Errors   VErrors
}

func (v *Validator) Required(p interface{}, k string) {
	failed := false
	switch v := p.(type) {
	case string:
		if v == "" {
			failed = true
		}
	case int:
		if v == 0 {
			failed = true
		}
	default:
		// @todo impl for other types
		failed = true
	}
	if failed {
		v.Errors.Set(k, VMsg{Msg: "valid_required", Params: map[string]interface{}{"Field": k}})
	}
}

func (v *Validator) AlphaDash(p, k string) {
	if reg_alpha_dash.MatchString(p) {
		v.Errors.Set(k, VMsg{Msg: "valid_alpha_dash", Params: map[string]interface{}{"Field": k}})
	}
}

func (v *Validator) Size(p, k string, min, max int) {
	if len(p) < min || len(p) > max {
		v.Errors.Set(k, VMsg{Msg: "valid_string_size", Params: map[string]interface{}{"Field": k, "Min": min, "Max": max}})
	}
}

func (v *Validator) Email(p, k string) {
	if !reg_email.MatchString(p) {
		v.Errors.Set(k, VMsg{Msg: "valid_email", Params: map[string]interface{}{"Field": k}})
	}
}

func (v *Validator) Range(p int, k string, min, max int) {
	if p < min || p > max {
		v.Errors.Set(k, VMsg{Msg: "valid_range", Params: map[string]interface{}{"Field": k, "Min": min, "Max": max}})
	}
}

func (v *Validator) NotContainStr(p, k string, ss []string) {
	for _, s := range ss {
		if p == s {
			v.Errors.Set(k, VMsg{Msg: "valid_not_contain", Params: map[string]interface{}{"Field": k, "Str": s}})
			break
		}
	}
}

// check in mongo collection if unique
func (v *Validator) UniqueDoc(k, c string, b bson.M) {
	sess := MgoSession.Copy()
	defer sess.Close()

	collection := sess.DB(MongoDbName).C(c)
	n, err := collection.Find(b).Select(bson.M{k: 1}).Count()
	if err == nil && n > 0 {
		v.Errors.Set(k, VMsg{Msg: "valid_unique", Params: map[string]interface{}{"Field": k}})
	}
}

func FormToBson(f map[string][]string) bson.M {

	b := make(bson.M)

	for k, v := range f {
		// exclude _xsrf field
		if v[0] != "" && k != "_xsrf" {
			b[k] = v[0]
		}
	}

	return b
}

// fmt field helper
func FmtString(prop string, actions []string) string {
	for _, v := range actions {
		switch v {
		case "ToLower":
			prop = strings.ToLower(prop)
		case "TrimSpace":
			prop = strings.TrimSpace(prop)
		case "Title":
			prop = strings.Title(prop)
		case "ToTitle":
			prop = strings.ToTitle(prop) // Unicode ToUpper
		case "ToUpper":
			prop = strings.ToUpper(prop)
		}
	}
	return prop
}

func genCacheKey(table string, method string, queries ...interface{}) string {
	var key string
	for _, v := range queries {
		key += ":" + fmt.Sprint(v)
	}
	key = strings.Replace(key, "} {", "},{", -1)
	key = strings.Replace(key, " ", ":", -1)

	return strings.ToLower(cachePrefix + table + ":" + method + key)
}

func cacheIsExist(key string) bool {
	return Cache.IsExist(key)
}

func cachePut(key string, data interface{}, timeout int64) error {
	// serialize only structs and bytes
	// prepare bytes buffer
	bCache := new(bytes.Buffer)
	encCache := gob.NewEncoder(bCache)
	err := encCache.Encode(data)
	Cache.Put(key, bCache, timeout)
	return err
}

func cacheGet(key string, data interface{}) error {
	decCache := gob.NewDecoder(bytes.NewBuffer(Cache.Get(key).([]byte)))
	err := decCache.Decode(data)
	return err
}

type VMsg struct {
	Msg    string
	Params map[string]interface{}
}

type BsonData struct {
	Raw bson.Raw
}

type VErrors map[string][]VMsg

// define model structs
type Geo struct {
	Type        string    `bson:"type"`
	Coordinates []float32 `bson:"coordsinates"`
}

type Address struct {
	City     string `bson:"city"`
	District string `bson:"dist"`
	Street   string `bson:"str"`
	RefLoc   string `bson:"ref"`
}

type NearStats struct {
	NScanned  uint32  `bson:"nscanned"`
	ObjLoaded uint32  `bson:"objectsLoaded"`
	AvrDis    float32 `bson:"avgDistance"`
	MaxDis    float32 `bson:"maxDistance"`
	time      int32   `bson:"time"`
}

// struct to store Near FoodServices result from mongo
type Near struct {
	Results bson.Raw
	Stats   NearStats
	Ok      float32
}

// validator func type
type VFunc func(val interface{}) (VMsg, bool)

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

// @todo embeded structs should be added bson and field dictionaries should
// create map of bson tag to field name
func BsonFieldDic(d interface {}) map[string]string {
	m := make(map[string]string)
	// should be pointer here
	s := reflect.ValueOf(d).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := typeOfT.Field(i)
		key := f.Tag.Get("bson")
		key = strings.Replace(key, " ", "", -1)
		key = strings.Split(key, ",")[0]
		m[key] = f.Name
	}
	return m
}
// create map of field name to bson tags
func FieldBsonDic(d interface {}) map[string]string {
	m := make(map[string]string)
	// should be pointer here
	s := reflect.ValueOf(d).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := typeOfT.Field(i)
		v := f.Tag.Get("bson")
		v = strings.Replace(v, " ", "", -1)
		v = strings.Split(v, ",")[0]
		m[f.Name] = v
	}
	return m

}

// @todo refactor to improve performance
func (v *VErrors) Set(key string, msg VMsg) {
	if msg.Msg == "" {
		return
	}
	vErrors := *v
	if len(vErrors) == 0 {
		vErrors = make(map[string][]VMsg)
	}
	vErrors[key] = append(vErrors[key], msg)
	*v = vErrors
}
func (v *VErrors) T(t i18n.TranslateFunc) map[string][]string {
	m := make(map[string][]string)
	if len(*v) == 0 {
		return m
	}
	for k, vm := range *v {
		for _, msg := range vm {
			// translate field names
			msg.Params["Field"] = t(msg.Params["Field"].(string))
			s := t(msg.Msg, msg.Params)
			m[k] = append(m[k], s)
		}
	}
	return m
}

// get bson field name from cached FieldDic, convenience func
func (u Address) Bson(f string) string {
	b, ok := FieldDic["Address"]["FieldBson"][f]
	if !ok {
		beego.Error("Address.Bson key in FieldDic does not exists " + f)
	}

	return b
}
