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
)

type BsonData struct {
	Raw bson.Raw
}

type DocModel interface {
	GetC() string
	GetIndex() []mgo.Index
	FmtFields()
	SetDefaults()
	Validate() (ValidationErrors, error)
	GetLocation() Geo
}

type ValidationErrors map[string][]string

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

func (v *ValidationErrors) Set(key, str string) {
	vErrors := *v
	if len(vErrors) == 0 {
		vErrors = make(map[string][]string)
	}
	vErrors[key] = append(vErrors[key], str)
	*v = vErrors
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

	// init indexes of models and panic if something wrong
	err = DocInitIndex(&FoodService{})
	panicOnErr(err)
	err = DocInitIndex(&User{})
	panicOnErr(err)
	// check new cache
	panicOnErr(errCache)

	// register Model structs for gob encoding
	gob.Register(User{})
	gob.Register(FoodService{})
	gob.Register(FacebookData{})

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

func DocCreate(m DocModel) (ValidationErrors, error) {
	// validate model before inserting
	vErrors, err := m.Validate()
	if err != nil || vErrors != nil {
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

func DocUpdate(q bson.M, m DocModel) (ValidationErrors, error) {
	// validate model before inserting
	vErrors, err := m.Validate()
	if err != nil || vErrors != nil {
		return vErrors, err
	}

	sess := MgoSession.Copy()
	defer sess.Close()

	// set defaults
	m.FmtFields()
	m.SetDefaults()

	collection := sess.DB(MongoDbName).C(m.GetC())
	err = collection.Update(q, m)

	return vErrors, err

}
