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