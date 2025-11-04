package main

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type transport struct {
	http.RoundTripper
}

// isLocalOrPrivate checks if a hostname resolves to localhost or private IP ranges
func isLocalOrPrivate(hostname string) bool {
	// Use net.SplitHostPort to properly handle IPv6 addresses
	host, _, err := net.SplitHostPort(hostname)
	if err != nil {
		// If there's no port, use the hostname as-is
		host = hostname
	}

	// Check for localhost
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return true
	}

	// Try to resolve the hostname to IP addresses
	ips, err := net.LookupIP(host)
	if err != nil {
		// If we can't resolve, block it to be safe
		return true
	}

	// Check if any resolved IP is private or loopback
	// This helps prevent DNS rebinding attacks by checking the actual resolved IPs
	for _, ip := range ips {
		if ip.IsLoopback() || ip.IsPrivate() {
			return true
		}
	}

	return false
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
		redirectURL, err := url.Parse(resp.Header.Get("Location"))
		if err == nil && redirectURL.Scheme == "https" {
			resp.Header.Set("Location", "http://"+req.Header.Get("X-Forwarded-Host")+"?url="+resp.Header.Get("Location"))
		}
	}
	return resp, nil
}

// handle a request send it to the server
func handleRequest(res http.ResponseWriter, req *http.Request) {
	// get the request URI
	server, err := url.Parse("https://" + host)
	if err != nil {
		http.Error(res, "Invalid server URL", http.StatusInternalServerError)
		return
	}

	reqURL, err := url.Parse(req.RequestURI)
	if err != nil {
		http.Error(res, "Invalid request URI", http.StatusBadRequest)
		return
	}

	if reqURL.Query().Get("url") != "" { // special proxy handling
		// extract the GET-Param url
		targetURL := reqURL.Query().Get("url")
		server, err = url.Parse(targetURL)
		if err != nil {
			http.Error(res, "Invalid target URL", http.StatusBadRequest)
			return
		}

		// Validate URL scheme to prevent SSRF attacks
		if server.Scheme != "https" && server.Scheme != "http" {
			http.Error(res, "Invalid URL scheme", http.StatusBadRequest)
			return
		}

		// Prevent access to localhost and private IP ranges
		if isLocalOrPrivate(server.Host) {
			http.Error(res, "Access to private addresses is not allowed", http.StatusForbidden)
			return
		}

		// replace request
		req.URL = server
		req.RequestURI = ""

		// mux host
		server, err = url.Parse(server.Scheme + "://" + server.Host)
		if err != nil {
			http.Error(res, "Invalid server URL", http.StatusInternalServerError)
			return
		}
	}

	// prepare reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(server)
	proxy.Transport = &transport{http.DefaultTransport}

	// update headers
	req.URL.Host = server.Host
	req.URL.Scheme = server.Scheme
	req.Header.Set("X-Forwarded-Host", req.Host)
	req.Host = server.Host

	if strings.HasPrefix(req.URL.Path, "/pins/") && req.URL.Host != host {
		res.Header().Add("Location", "http://"+req.Header.Get("X-Forwarded-Host")+"?url="+"https://"+host+req.URL.Path)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	// run the proxy
	proxy.ServeHTTP(res, req)
}

func runProxy(disableSideloading bool) {
	// handle simple information path
	http.HandleFunc("/info", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = res.Write([]byte("The Plex proxy service is running on " + req.Host))
	})

	//handle widgetlist for sideloading
	http.HandleFunc("/widgetlist.xml", func(res http.ResponseWriter, req *http.Request) {
		buf, err := retreiveZipFile()
		if err != nil {
			res.Header().Set("Content-Type", "text/plain; charset=utf-8")
			res.WriteHeader(http.StatusInternalServerError)
			_, _ = res.Write([]byte(err.Error()))
			log.Println(err)
			return
		}
		xml := `<?xml version="1.0" encoding="UTF-8"?>
		<rsp stat="ok">
				<list>
						<widget id="` + modifiedAppName + `">
								<title>Plex</title>
								<compression size="` + strconv.Itoa(len(buf)) + `" type="zip"/>
								<description/>
								<download>http://` + req.Host + `/` + modifiedAppFile + `</download>
						</widget>
				</list>
		</rsp>`

		res.Header().Set("Content-Type", "application/xml; charset=utf-8")
		_, _ = res.Write([]byte(xml))
	})

	// handle app-deployment
	http.HandleFunc("/"+modifiedAppFile, func(res http.ResponseWriter, req *http.Request) {
		buf, err := retreiveZipFile()
		if err != nil {
			res.Header().Set("Content-Type", "text/plain; charset=utf-8")
			res.WriteHeader(http.StatusInternalServerError)
			_, _ = res.Write([]byte(err.Error()))
			log.Println(err)
			return
		}
		// write the http-response
		res.Header().Set("Content-Type", "application/zip")
		_, _ = res.Write(buf)
	})

	// start real proxy
	http.HandleFunc("/", handleRequest)

	// Configure HTTP server with timeouts to prevent DoS attacks
	serverMain := &http.Server{
		Addr:              ":" + port,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// try to handle everything on port 80 aswell for serving the app
	// Note: this will not work on non-rooted android because only high-ports can be used
	// The goroutine will exit after logging the error, which is expected behavior
	// for cases where port 80 requires elevated privileges
	go func() {
		if !disableSideloading {
			log.Println("Trying to start app-deployer on port 80 ...")
			server80 := &http.Server{
				Addr:              ":80",
				ReadTimeout:       15 * time.Second,
				WriteTimeout:      15 * time.Second,
				IdleTimeout:       60 * time.Second,
				ReadHeaderTimeout: 5 * time.Second,
			}
			if err := server80.ListenAndServe(); err != nil {
				log.Printf("Port 80 server error (this is expected on non-rooted systems): %v\n", err)
			}
		}
	}()

	log.Println("Server starting on Port " + port + " ...")
	log.Fatal(serverMain.ListenAndServe())
}
