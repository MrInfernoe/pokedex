package pokecache_test

import (
	// "testing"
	// "pokedex/internal/pokecache"
)

// func TestAddGet(t *testing.T) {
// 	cases := []struct{
// 		inputKey 	string
// 		inputBytes 	[]byte
// 		expected	*pokecache.Cache
// 	}{
// 		{
// 			inputKey: "this is a string",
// 			inputBytes: []byte("so this is data converted to bytes"),
// 			expected: pokecache.Cache{},
// 		},
// 	}
// }

/* start from boot.dev
func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}
end from boot.dev */



// /*
// Check whether correct Cache is made
// Check whether cacheEntry exists during duration
// Check if cacheEntry access refreshes duration
// Check whether cacheEntry deleted after duration
// */

// func TestNewCache(t *testing.T){
// 	cases := []struct {
// 		inputInterval time.Duration
// 		expectedCache Cache
// 	}{

// 	}
// }