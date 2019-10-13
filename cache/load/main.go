package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var (
	addr     = flag.String("addr", "localhost:8080", "'cache' server address to load")
	sdur     = flag.String("dur", "10s", "how long")
	verbose  = flag.Bool("v", false, "print retrieved cache values and error strings")
	vverbose = flag.Bool("vv", false, "very verbose")
)

var (
	m       sync.Map
	keys    []string
	letters = "abcdefghijklmnopqrstuvwxyz"
)

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randomAdd() *http.Request {
	var k, v string
	rnd := rand.Intn(100)
	if rnd < 32 {
		k = keys[rnd%len(keys)]
		iv, ok := m.Load(k)
		if !ok {
			panic(fmt.Sprintf("m.Load(%s)", k))
		}
		v = iv.(string)
	} else {
		ch := randomString(1)
		k = "key-" + ch
		v = "value-" + ch // don't care
	}
	q := fmt.Sprintf("http://%s/add?k=%s&v=%v", *addr, k, v)
	r, err := http.NewRequest("GET", q, nil)
	if err != nil {
		panic(err)
	}
	return r
}

func randomGet() *http.Request {
	var k string
	rnd := rand.Intn(100)
	if rnd < 32 {
		k = keys[rnd%len(keys)]
	} else {
		k = "key-" + randomString(1)
	}
	q := fmt.Sprintf("http://%s/get?k=%s", *addr, k)
	r, err := http.NewRequest("GET", q, nil)
	if err != nil {
		panic(err)
	}
	return r
}

func randomRequest() *http.Request {
	rnd := rand.Intn(100)
	if rnd < 50 {
		return randomGet()
	}
	return randomAdd()
}

func setup() time.Duration {
	flag.Parse()
	if *vverbose {
		*verbose = true
	}

	dur, err := time.ParseDuration(*sdur)
	if err != nil {
		panic(fmt.Sprint("dur:", err))
	}

	rand.Seed(time.Now().UnixNano())

	m.Store("hello", "golab")
	m.Store("language", "Go")
	m.Store("version", "1.13")
	m.Store("topic", "monitoring")
	m.Store("pi", fmt.Sprint(math.Pi))
	m.Store("tsdb", "prometheus")
	m.Store("year", "2019")
	m.Store("month", "october")

	m.Range(func(k, v interface{}) bool {
		keys = append(keys, k.(string))
		return true
	})
	if len(keys) != 8 {
		panic("should remain 8 or bad stuff will happen!")
	}
	return dur
}

func main() {
	dur := setup()

	ncpu := runtime.NumCPU()
	totalRequests := make([]int, ncpu)
	totalErrors := make([]int, ncpu)

	wg := sync.WaitGroup{}
	wg.Add(ncpu)
	for p := 0; p < ncpu; p++ {
		p := p
		go func() {
			defer wg.Done()
			deadline := time.NewTimer(dur)
			for {
				timer := time.NewTicker(time.Duration((rand.Intn(20) + 50)) * time.Millisecond)
				select {
				case <-deadline.C:
					return
				case <-timer.C:
				}
				req := randomRequest()
				if *vverbose {
					fmt.Println("request:", req.URL.String())
				}
				resp, err := http.DefaultClient.Do(req)
				totalRequests[p]++
				if err != nil {
					totalErrors[p]++
					if *verbose {
						fmt.Println("request error:", err)
					}
					continue
				}
				buf, _ := ioutil.ReadAll(resp.Body)
				if *verbose && len(buf) != 0 {
					fmt.Println(req.URL.String(), "responded", string(buf))
				}
				resp.Body.Close()
			}
		}()
	}
	wg.Wait()

	nrequests, nerrors := 0, 0
	for _, v := range totalErrors {
		nerrors += v
	}
	for _, v := range totalRequests {
		nrequests += v
	}

	fmt.Println("-----------")
	fmt.Printf("requests %v\n", nrequests)
	fmt.Printf("errors   %v\t%.01f%%\n", nerrors, float64(nerrors)/float64(nrequests))
}
