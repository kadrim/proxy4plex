package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
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
	server, _ := url.Parse("https://" + host)
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
		res.Write([]byte("The Plex proxy service is running on " + req.Host))
	})

	//handle widgetlist for sideloading
	http.HandleFunc("/widgetlist.xml", func(res http.ResponseWriter, req *http.Request) {
		buf, err := retreiveZipFile()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
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

		res.Write([]byte(xml))
	})

	// handle app-deployment
	http.HandleFunc("/"+modifiedAppFile, func(res http.ResponseWriter, req *http.Request) {
		buf, err := retreiveZipFile()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			log.Println(err)
			return
		}
		// write the http-response
		res.Write(buf)
	})

	// start real proxy
	http.HandleFunc("/", handleRequest)

	// try to handle everything on port 80 aswell for serving the app
	// Note: this will not work on non-rooted android because only high-ports can be used
	go func() {
		if !disableSideloading {
			log.Println("Trying to start app-deployer on port 80 ...")
			http.ListenAndServe(":80", nil)
		}
	}()

	log.Println("Server starting on Port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
