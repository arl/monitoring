# Cache

```
$ ./cache -h
Usage of ./cache:
  -addr string
        server listen address (default ":8080")
  -size int
        LRU cache size
```

`cache` is a HTTP server that wraps a very basic LRU cache (Least Recently Used) cache.


## Add a value to the cache

 - Use the `/add` endpoint and provide a KEY and a VALUE
 -  http://host:port/add?k=KEY&v=VALUE

## Get a value from the cache

 - Use the `/get` endpoint and provide a KEY
 -  http://host:port/get?k=KEY

The response body will be either empty, indicating a cache miss, or the cached value.


## Use with Prometheus

Once the server is instrumented, Prometheus needs to periodically scrape a new /metrics endpoint.

### Binary installation of Prometheus

Use the provided `prometheus.yml` and replace the `host:port` in the targets list with 
`127.0.0.1:8080`.  
Access the Prometheus web-ui at 127.0.0.1:9090 (default).  
After hacking the Go code in the `cache` package, just `go run` or `go build` it.

### Use Prometheus/Grafana stack with docker-compose

Your current directory should be the directory containing this README.
Run `docker-compose up --build`.
After hacking the Go code in the `cache` package:
 - open a new terminal (keep docker-compose running)
 - go build .
 - docker-compose restart app