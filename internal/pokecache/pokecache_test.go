package pokecache_test

import (
	"testing"
	"pokedex/internal/pokecache"
	"time"
)

func TestAddGet(t *testing.T) {
	interval := 5*time.Second
	cases := []struct{
		name	string
		key		string
		val		[]byte
	}{
		{
			name: "basic1",
			key: "https://fu.bar",
			val: []byte("I am data"),
		},
		{
			name: "basic2",
			key: "https://bar.fu",
			val: []byte("I am other data, not that data"),
		},
	}

	for _, entry := range cases {
		runOk := t.Run(entry.name, func(t *testing.T) {
				cache := pokecache.NewCache(interval)
				cache.Add(entry.key, entry.val)
				val, exists := cache.Get(entry.key)
				if !exists {
					t.Errorf("key not found")
				}
				for i, valByte := range val {
					if valByte != entry.val[i] {
						t.Errorf("not all values same inside cache")
					}
				}
			})
		if !runOk {
			t.Errorf("could not run")
		}
	}
}

func TestReaping(t *testing.T) {
	const interval = 100*time.Millisecond
	cache := pokecache.NewCache(interval)
	cache.Add("https://fu.bar", []byte("this'll be data"))

	// checking reaping
	if _, exists := cache.Get("https://fu.bar"); !exists {
		t.Errorf("expected to find data")
	}

	time.Sleep(200*time.Millisecond)

	if _, exists := cache.Get("https://fu.bar"); exists {
		t.Errorf("expected to NOT find data")
	}


	// checking duration reset after get
	cache.Add("https://fu.bar", []byte("this'll be data"))

	if _, exists := cache.Get("https://fu.bar"); !exists {
		t.Errorf("expected to find data")
	}
	// reseting duration with Get
	time.Sleep(70*time.Millisecond)
	if _, exists := cache.Get("https://fu.bar"); !exists {
		t.Errorf("expected to find data")
	}
	time.Sleep(70*time.Millisecond)
	if _, exists := cache.Get("https://fu.bar"); !exists {
		t.Errorf("expected to find data")
	}
	time.Sleep(70*time.Millisecond)
	if _, exists := cache.Get("https://fu.bar"); !exists {
		t.Errorf("expected to find data")
	}
	// 210 Milliseconds passed since added

	time.Sleep(200*time.Millisecond)

	if _, exists := cache.Get("https://fu.bar"); exists {
		t.Errorf("expected to NOT find data")
	}
}

/*
make interval
makes cases: string, []byte
for range of cases
	run anonymous function of
		new cache
		adding entry
		getting entry
		given enough time check if entry still exists
		given time less than duration get entry (to reset createdAt) 
			then check if still exists after standard duration
			then check is deleted after duration reset
*/

/*
make cache
add entry
check it exists
sleep for time then check again
*/

/* start from boot.dev
func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
end from boot.dev */