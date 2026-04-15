package pokecache_test

import (
	"testing"
	"pokedex/internal/pokecache"
)

/*
Check whether correct Cache is made
Check whether cacheEntry exists during duration
Check if cacheEntry access refreshes duration
Check whether cacheEntry deleted after duration
*/

func TestNewCache(t *testing.T){
	cases := []struct {
		inputInterval time.Duration
		expectedCache Cache
	}{

	}
}