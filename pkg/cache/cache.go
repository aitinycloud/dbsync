package cache

import (
	"fmt"
	"time"

	gcache "github.com/patrickmn/go-cache"
)

//GCache Loacl GO cacha from GCache
var GCache *gcache.Cache

var cacheFlag bool = false

//CacheInit Cache Init
func CacheInit() bool {
	GCache = gcache.New(gcache.NoExpiration, time.Minute)
	tmpkey := "tmpkey1"
	tmpvalue := "tmpvalue1"
	GCache.Set(tmpkey, tmpvalue, gcache.NoExpiration)
	if val, found := GCache.Get("tmpkey1"); found {
		if val != tmpvalue {
			panic(fmt.Sprintln("CacheInit GCache.Get fail."))
		}
	} else {
		panic(fmt.Sprintln("CacheInit GCache.Get fail. no found"))
	}
	cacheFlag = true
	return true
}

func IsCacheInit() bool {
	return cacheFlag
}
