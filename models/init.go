package models

import (
	"bytes"
	"fmt"
	"strings"

	"encoding/gob"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"regexp"
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
)

// convenience map for i18n translation parameters map
type tm map[string]interface{}

type VMsg struct {
	Msg   string
	Param tm
}

type BsonData struct {
	Raw bson.Raw
}

type DocModel interface {
	GetC() string
	GetIndex() []mgo.Index
	FmtFields()
	SetDefaults()
	Validate(bs bson.M) VErrors
	GetLocation() Geo
}

type VErrors map[string][]VMsg

// define model structs
type Geo struct {
	Type        string    `bson:"type"`
	Coordinates []float32 `bson:"coordinates"`
}

type Address struct {
	City     string `bson:"city"`
	District string `bson:"district"`
	Street   string `bson:"street"`
	RefLoc   string `bson:"ref_loc"`
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
func BsonFieldDic(d DocModel) map[string]string {
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

func FieldBsonDic(d DocModel) map[string]string {
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

func AndSet(val interface{}, tfs []VFunc) []VMsg {
	var msgs []VMsg
	// first fn should be Required func
	// if it return true continue validation
	if len(tfs) == 0 {
		return msgs
	}
	_, ok := tfs[0](val)
	if !ok {
		return msgs
	}

	for i := 1; i < len(tfs); i++ {
		msg, ok := tfs[i](val)
		if !ok {
			msgs = append(msgs, msg)
		}
	}
	return msgs
}

func (v *VErrors) Set(key string, ss []VMsg) {
	if len(ss) == 0 {
		return
	}
	vErrors := *v
	if len(vErrors) == 0 {
		vErrors = make(map[string][]VMsg)
	}
	vErrors[key] = ss
	*v = vErrors
}

// validation functions
func Required(b bool) VFunc {
	return func(val interface{}) (VMsg, bool) {
		var vm VMsg
		if b == true {
			return vm, true
		} else {
			isR := false
			switch val.(type) {
			case string:
				if val.(string) != "" {
					isR = true
				}
			case int:
				if val.(int) != 0 {
					isR = true
				}
			}
			if isR {
				return vm, true
			} else {
				return vm, false
			}
		}
	}
}

func MinInt(min int) VFunc {
	return func(val interface{}) (VMsg, bool) {
		var vm VMsg
		if val.(int) >= min {
			return vm, true
		}
		vm.Msg = "valid_min"
		vm.Param = tm{"Min": min}
		return vm, false
	}
}

func RangeStr(min, max int) VFunc {
	return func(val interface{}) (VMsg, bool) {
		var vm VMsg
		l := len(val.(string))
		if l >= min && l <= max {
			return vm, true
		}
		vm.Msg = "valid_range"
		vm.Param = tm{"Min": min, "Max": max}
		return vm, false
	}
}

func NotEmptyStr() VFunc {
	return func(val interface{}) (VMsg, bool) {
		var vm VMsg
		if val.(string) != "" {
			return vm, true
		}
		vm.Msg = "valid_not_empty"
		return vm, false
	}
}

func InSetStr(ss []string) VFunc {
	return func(val interface{}) (VMsg, bool) {
		var vm VMsg
		for _, v := range ss {
			if val.(string) == v {
				return vm, true
			}
		}
		vm.Msg = "valid_set"
		vm.Param = tm{"Set": strings.Join(ss, ", ")}
		return vm, false
	}
}

func ValidEmail() VFunc {
	emailPattern := regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")
	return func(val interface{}) (VMsg, bool) {
		var vm VMsg
		if emailPattern.MatchString(val.(string)) {
			return vm, true
		}
		vm.Msg = "valid_email"
		return vm, false
	}
}

func NotContainStr(ss []string) VFunc {
	return func(val interface{}) (VMsg, bool) {
		var vm VMsg
		for _, v := range ss {
			if strings.Index(val.(string), v) != -1 {
				vm.Msg = "valid_not_contain"
				vm.Param = tm{"Set": strings.Join(ss, ", ")}
				return vm, false
			}
		}
		return vm, true
	}
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
	var foodService FoodService
	FieldDic = make(map[string]map[string]map[string]string)
	FieldDic["User"] = make(map[string]map[string]string)
	// define maps bsonToUser and userToBson field names
	FieldDic["User"]["FieldBson"] = FieldBsonDic(&user)
	FieldDic["User"]["BsonField"] = BsonFieldDic(&user)

	// init indexes of models and panic if something wrong
	err = DocInitIndex(&foodService)
	panicOnErr(err)
	err = DocInitIndex(&user)
	panicOnErr(err)
	// check new cache
	panicOnErr(errCache)

	// register Model structs for gob encoding
	gob.Register(user)
	gob.Register(foodService)
	gob.Register(FacebookData{})
	gob.Register(GoogleData{})

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
	vErrors := m.Validate(bson.M{})
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
	vErrors := m.Validate(flds)
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
