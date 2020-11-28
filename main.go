package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	port   = "3000"
	server = "http://plex.tv"
)

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// check if result is a redirect and handle that accordingly
	switch resp.StatusCode {
	case 301:
		fallthrough
	case 302:
		redirectURL, _ := url.Parse(resp.Header.Get("Location"))
		if redirectURL.Scheme == "https" {
			resp.Header.Set("Location", "http://"+req.Header.Get("X-Forwarded-Host")+"?url="+resp.Header.Get("Location"))
		}
	}
	return resp, nil
}

// handle a request send it to the server
func handleRequest(res http.ResponseWriter, req *http.Request) {
	// get the request URI
	server, _ := url.Parse(server)
	reqURL, _ := url.Parse(req.RequestURI)
	if reqURL.Query().Get("url") != "" { // special proxy handling
		// extract the GET-Param url
		server, _ = url.Parse(reqURL.Query().Get("url"))

		// replace request
		req.URL = server
		req.RequestURI = ""

		// mux host
		server, _ = url.Parse(server.Scheme + "://" + server.Host)
	}

	// prepare reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(server)
	proxy.Transport = &transport{http.DefaultTransport}

	// update headers
	req.URL.Host = server.Host
	req.URL.Scheme = server.Scheme
	req.Header.Set("X-Forwarded-Host", req.Host)
	req.Host = server.Host

	// run the proxy
	proxy.ServeHTTP(res, req)
}

func main() {
	// handle simple information path
	http.HandleFunc("/info", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("The Plex proxy service is running"))
	})

	// start real server
	http.HandleFunc("/", handleRequest)

	log.Println("Server starting on Port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
