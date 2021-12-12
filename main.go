package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const (
	port              = "3000"
	host              = "plex.tv"
	officialAppURL    = "https://www.dropbox.com/s/f17hx2w7tvofjqr/Plex_2.014_11112020.zip?dl=1"
	officialAppChksum = "8c6b2bb25a4c2492fd5dbde885946dcb6b781ba292e5038239559fd7a20e707e"
	modifiedAppName   = "Plex_2.014_net"
	modifiedAppFile   = "Plex_2.014_11112020_net.zip"
)

var modifiedBuffer []byte

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

func extractZipFile(zipFile *zip.File) ([]byte, error) {
	file, err := zipFile.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func retreiveZipFile() ([]byte, error) {
	// only retreive data once
	if modifiedBuffer != nil {
		return modifiedBuffer, nil
	}

	// download the original zip
	resp, err := http.Get(officialAppURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// checksum
	chksum := sha256.Sum256(body)
	if hex.EncodeToString(chksum[:]) == officialAppChksum {
		log.Println("checksum match. going on ...")
	} else {
		err := errors.New("checksum did not match, aborting")
		return nil, err
	}

	// attach the zipWriter to the buffer
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// create zipReader
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}

	// regexp pattern to remove first directory
	regex := regexp.MustCompile(`^Plex/`)

	// process all files, read each file and write it to new zip
	for _, zipFile := range zipReader.File {
		fmt.Println("Processing file:", zipFile.Name)
		newFileName := regex.ReplaceAllString(zipFile.Name, "")
		fmt.Println("Replacement:", newFileName)
		bytes, err := extractZipFile(zipFile)
		if err != nil {
			return nil, err
		}

		// write file to new zip
		file, err := w.Create(newFileName)
		if err != nil {
			return nil, err
		}
		_, err = file.Write(bytes)
		if err != nil {
			return nil, err
		}
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	// write data to global var
	modifiedBuffer = buf.Bytes()

	return modifiedBuffer, nil
}

func main() {
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
		log.Println("Trying to start app-deployer on port 80 ...")
		http.ListenAndServe(":80", nil)
	}()

	log.Println("Server starting on Port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
