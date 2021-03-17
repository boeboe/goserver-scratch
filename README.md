## goserver scratch container

A small scratch based docker container with and GoLang HTTP server and client, compressed with upx

This container can be used to generate a network of microservices for performance testing

https://github.com/boeboe/goserver-scratch

https://hub.docker.com/repository/docker/boeboe/goserver-scratch

### Usage

#### Standalone

```
start a go http server and client with command line flag or environment variables [ENV]

Usage of goserver:
  -b3_trace
    	[B3_TRACE] Enable B3 header propagation for traces (default disabled)
  -help
    	Print this help
  -request_size int
    	[REQUEST_SIZE] Request body size (default 50)
  -response_size int
    	[RESPONSE_SIZE] Respose body size (default 50)
  -server_name string
    	[SERVER_NAME] Name of server or hostname (defaults to hostname)
  -server_port int
    	[SERVER_PORT] Server port (default 8080)
  -upstream_host string
    	[UPSTREAM_HOST] Name of upstream server
  -upstream_port int
    	[UPSTREAM_PORT] Upstream server Port (default 8080)
  -verbose
    	[VERBOSE] Verbose output (default disabled)
```

#### Docker

```
docker run -i -t --rm -p 8080:8080 --name="goserver-scratch" boeboe/goserver-scratch:1.0.0 -help
```
