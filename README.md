# http-client-test

The purpose of this project to help me understand [Go http client](https://golang.org/pkg/net/http/) configuration and its effects on performance.

It contains a single webapp that calls another service and return its response.

The other service is [fake-service](https://github.com/johnmuth/fake-service), which returns a hard-coded response after a random delay.

## http-client-test endpoints

- /api
    - response: `{"requestid":"6ba7b810-9dad-11d1-80b4-00c04fd430c8","qux":"flubber"}`
        - requestid is a unique id to help correlate events in the logs
        - qux is from the response from fake-service. 
        
- /internal/healtcheck
    - for load balancer
    
## net/http Client

To send HTTP requests from within your Go code, the standard way is to use Go's [net/http package](https://golang.org/pkg/net/http/).

It provides convenient methods for sending requests:

```go
resp1, err := http.Get("http://example.com/")
resp2, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
resp3, err := http.PostForm("http://example.com/form",url.Values{"key": {"Value"}, "id": {"123"}})
```

When you use those methods you're using a default `http.Client` with default values for a bunch of options that you might want to override:

```go
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 2,
			DialContext: (&net.Dialer{
        			Timeout:   30 * time.Second,
		        	KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 0,
	}
```

Those defaults are okay to get started, but you'll definitely want to override at least some of them to use in production.

For example,

* **Timeout: 0** : AKA no timeout: outgoing request can hang forever
* **MaxIdleConnsPerHost: 2** : If you're doing a lot of requests to the same host, you probably want to allow more idle connections per host.

To make it easy to experiment, all of the options are configurable in this project via environment variables. Look at [docker-compose.yml](docker-compose.yml) to see them all. Read the [net/http package](https://golang.org/pkg/net/http/) source to understand what they all mean. (I'll also add words here to summarise what I learn.)

## httptrace

The [httptrace package](https://golang.org/pkg/net/http/httptrace) provides a nice way to add logging and/or metrics to events within HTTP client requests.

[service.go](service.go), shows how to add httptrace to an existing http client request/response flow.

The result is a lot of log messages tracing the life of a single request, starting with "get connection" and ending with "put idle connection" - returning the connection to the connection pool:

```bash
{"level":"info","msg":"About to get connection","requestid":"0519190b-0bb6-4618-a974-7492776b40d9","time":"2017-09-03T13:13:26.283405633Z"}
{"idletime":6797540509,"level":"info","msg":"Got connection","requestid":"0519190b-0bb6-4618-a974-7492776b40d9","reused":true,"time":"2017-09-03T13:13:26.283481947Z","wasidle":true}
{"level":"info","msg":"Wrote headers","requestid":"0519190b-0bb6-4618-a974-7492776b40d9","time":"2017-09-03T13:13:26.28354208Z"}
{"level":"info","msg":"Wrote request","requestid":"0519190b-0bb6-4618-a974-7492776b40d9","time":"2017-09-03T13:13:26.28358753Z"}
{"level":"info","msg":"First response byte!","requestid":"0519190b-0bb6-4618-a974-7492776b40d9","time":"2017-09-03T13:13:26.284466249Z"}
{"err":null,"level":"info","msg":"Put idle connection","requestid":"0519190b-0bb6-4618-a974-7492776b40d9","time":"2017-09-03T13:13:26.285229498Z"}
```

## requestid

To make sense of the detailed log messages when the application is handling lots of requests concurrently, it generates a unique **requestid** in [handler.go](handler.go) and uses it throughout.


## locust

To do the load testing I'm using [locust](http://locust.io), because it looks nice, I like Python, and I'm bored of [Gatling](http://gatling.io) (maybe just bored of apologising for Scala!)




