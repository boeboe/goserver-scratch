package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"

	"github.com/boeboe/flag"
)

type ServerConfig struct {
	b3trace      bool
	json         bool
	requestsize  int
	responsesize int
	servername   string
	serverport   int
	upstreamhost string
	upstreamport int
	verbose      bool
}

// B3 headers documented by Istio (https://istio.io/latest/faq/distributed-tracing)
var traceHeaders = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"b3",
	"x-ot-span-context",
	"x-cloud-trace-context",
	"traceparent",
	"grpc-trace-bin"}

func (sc ServerConfig) get(res http.ResponseWriter, req *http.Request) {
	if sc.verbose {
		log.Printf("========== Incoming request received ==========")
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Println(err)
		}
		log.Println(string(requestDump))
	}

	if sc.upstreamhost != "" {
		upstreamreq, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/",
			sc.upstreamhost, sc.upstreamport), bytes.NewBufferString(sc.generateHttpBody()))

		if err != nil {
			log.Fatal(err)
		}

		if sc.b3trace {
			upstreamreq = copyTraceHeaders(req, upstreamreq)
		}

		client := &http.Client{}

		upstreamresp, errdo := client.Do(upstreamreq)

		for errdo != nil {
			if sc.verbose {
				log.Printf("Upstream request failed, sleep 1s and try again")
			}
			time.Sleep(1 * time.Second)
			upstreamresp, errdo = client.Do(upstreamreq)
		}

		defer upstreamresp.Body.Close()

		if sc.verbose {
			if sc.verbose {
				log.Printf("========== Incoming response received ==========")
				requestDump, err := httputil.DumpResponse(upstreamresp, true)
				if err != nil {
					fmt.Println(err)
				}
				log.Println(string(requestDump))
			}
		}
	}

	fmt.Fprint(res, sc.generateHttpBody())
}

func (sc ServerConfig) generateHttpBody() string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, sc.responsesize)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}

	if sc.json {
		return fmt.Sprintf("{\"timestamp\": \"%s\" ,\"server_name\": \"%s\", \"response\": \"%s\"}\n",
			time.Now().Format("2006-01-02 15:04:05.000000"), sc.servername, string(b))
	} else {
		return string(b)
	}
}

func copyTraceHeaders(req *http.Request, upstreamreq *http.Request) *http.Request {
	headers := req.Header

	for _, traceHeader := range traceHeaders {
		if headers.Get(traceHeader) != "" {
			upstreamreq.Header.Set(traceHeader, headers.Get(traceHeader))
		}
	}

	return upstreamreq
}

func main() {

	hostname, _ := os.Hostname()

	b3trace := flag.Bool("b3_trace", false, "[B3_TRACE] Enable B3 header propagation for traces (default disabled)")
	help := flag.Bool("help", false, "Print this help")
	json := flag.Bool("json", false, "[JSON] Enable JSON for request/response bodies")
	requestsize := flag.Int("request_size", 50, "[REQUEST_SIZE] Request body size")
	responsesize := flag.Int("response_size", 50, "[RESPONSE_SIZE] Respose body size")
	servername := flag.String("server_name", hostname, "[SERVER_NAME] Name of server or hostname")
	serverport := flag.Int("server_port", 8080, "[SERVER_PORT] Server port")
	upstreamhost := flag.String("upstream_host", "", "[UPSTREAM_HOST] Name of upstream server")
	upstreamport := flag.Int("upstream_port", 8080, "[UPSTREAM_PORT] Upstream server Port")
	verbose := flag.Bool("verbose", false, "[VERBOSE] Verbose output (default disabled)")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	sc := ServerConfig{
		b3trace:      *b3trace,
		json:         *json,
		requestsize:  *requestsize,
		responsesize: *responsesize,
		servername:   *servername,
		serverport:   *serverport,
		upstreamhost: *upstreamhost,
		upstreamport: *upstreamport,
		verbose:      *verbose}

	log.Printf("Bootstrap configuration:\n\tb3_trace: %t\n\tjson: %t\n\trequest_size: %d\n\t"+
		"response_size: %d\n\tserver_name: %s\n\tserver_port: %d\n\tupstream_host: %s\n\t"+
		"upstream_port: %d\n\tverbose: %t\n", sc.b3trace, sc.json, sc.requestsize, sc.responsesize,
		sc.servername, sc.serverport, sc.upstreamhost, sc.upstreamport, sc.verbose)

	http.HandleFunc("/", sc.get)

	log.Printf("Going to listen on port %d\n", *serverport)

	if err := http.ListenAndServe(":"+strconv.Itoa(*serverport), nil); err != nil {
		log.Fatal(err)
	}
}
