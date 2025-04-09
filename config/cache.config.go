package config

import (
	"time"

	"github.com/viccon/sturdyc"
)

var CacheClient *sturdyc.Client[[]byte]

func ConfigCacheInit() {
	capacity := 10000
	numShards := 10
	ttl := 1 * time.Minute
	evictionPercentage := 10
	CacheClient = sturdyc.New[[]byte](capacity, numShards, ttl, evictionPercentage)
}
