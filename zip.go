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
	"os"
	"regexp"
)

var modifiedBuffer []byte

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
		log.Println("reusing already prepared net-sideloader")
		return modifiedBuffer, nil
	}

	var zipData []byte
	var download bool

	if _, err := os.Stat(officialAppFile); err == nil { // check if the plex-app has already been downloaded
		log.Println("plex-app exists in local directory, not downloading it again")
		zipData, err = ioutil.ReadFile(officialAppFile)
		if err != nil {
			return nil, err
		}
	} else {
		// download the original zip
		log.Println("downloading application from " + officialAppURL + " ... ")
		resp, err := http.Get(officialAppURL)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		log.Println("done!")

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		} else {
			zipData = body
		}

		download = true
	}

	// checksum
	chksum := sha256.Sum256(zipData)
	if hex.EncodeToString(chksum[:]) == officialAppChksum {
		log.Println("checksum match. going on ...")
	} else {
		var message string
		message = "checksum did not match, aborting"
		if !download {
			message = message + "\nPlease delete existing file " + officialAppFile + " to force redownload"
		}
		err := errors.New(message)
		log.Println(message)

		return nil, err
	}

	if download {
		// write zipData to local file for caching
		err := ioutil.WriteFile(officialAppFile, zipData, 0664)
		if err != nil {
			log.Fatal("could not save downloaded file, going on anyway")
		}
	}

	// attach the zipWriter to the buffer
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// create zipReader
	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
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

	log.Println("app ready for net-sideloading!")

	return modifiedBuffer, nil
}
