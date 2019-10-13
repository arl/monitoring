package main

import (
	"fmt"
	"testing"
)

func TestLRUCache(t *testing.T) {
	c := NewLRUCache(2)
	c.Add("hello", 1)
	c.Add("golab", 2)
	c.Add("2019", 3)

	var (
		value interface{}
		ok    bool
	)
	value, ok = c.Get("golab")
	if !ok || value != 2 {
		t.Fatalf(`c["hello"] = (%v %t), want (%v, %v)`, value, ok, 2, true)
	}

	value, ok = c.Get("2019")
	if !ok || value != 3 {
		t.Fatalf(`c["2019"] = (%v %t), want (%v, %v)`, value, ok, 3, true)
	}

	value, ok = c.Get("hello")
	if ok {
		t.Fatalf(`c["hello"] = (%v %t), want (%v, %v)`, value, ok, nil, false)
	}

	c.Add("golab", 6)
	c.Add("key", 2)
	c.Add("key", 3)
	c.Add("key", 4)
	c.Add("key", 5)
	c.Add("key", 6)
	c.Add("key", 7)

	value, ok = c.Get("golab")
	if !ok || value != 6 {
		t.Fatalf(`c["golab"] = (%v %t), want (%v, %v)`, value, ok, 6, true)
	}

	value, ok = c.Get("key")
	if !ok || value != 7 {
		t.Fatalf(`c["key"] = (%v %t), want (%v, %v)`, value, ok, 7, true)
	}
}

func (c *LRUCache) fill() {
	for i := 0; i < c.maxcap; i++ {
		c.Add(fmt.Sprintf("k-%d", i), nil)
	}
}

var sink interface{}

func benchmarkLRUCacheHit(b *testing.B, maxcap int) {
	c := NewLRUCache(maxcap)
	c.fill()

	for n := 0; n < b.N; n++ {
		sink, _ = c.Get("hit")
	}
}

func benchmarkLRUCacheMiss(b *testing.B, maxcap int) {
	c := NewLRUCache(maxcap)
	c.fill()

	for n := 0; n < b.N; n++ {
		sink, _ = c.Get("miss")
	}
}

func BenchmarkLRUCacheHit_100(b *testing.B)    { benchmarkLRUCacheHit(b, 100) }
func BenchmarkLRUCacheHit_1000(b *testing.B)   { benchmarkLRUCacheHit(b, 1000) }
func BenchmarkLRUCacheHit_10000(b *testing.B)  { benchmarkLRUCacheHit(b, 10000) }
func BenchmarkLRUCacheHit_100000(b *testing.B) { benchmarkLRUCacheHit(b, 100000) }

func BenchmarkLRUCacheMiss_100(b *testing.B)    { benchmarkLRUCacheMiss(b, 100) }
func BenchmarkLRUCacheMiss_1000(b *testing.B)   { benchmarkLRUCacheMiss(b, 1000) }
func BenchmarkLRUCacheMiss_10000(b *testing.B)  { benchmarkLRUCacheMiss(b, 10000) }
func BenchmarkLRUCacheMiss_100000(b *testing.B) { benchmarkLRUCacheMiss(b, 100000) }
