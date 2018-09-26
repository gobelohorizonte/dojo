/*
* Go Library (C) 2017 Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
* @project     FileSerever
* @package     main
* @author      @jeffotoni
* @size        23/09/2018
*
 */

package handler

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

import (
	"github.com/jeffotoni/fileserver/pkg/cryptf"
	"github.com/jeffotoni/fileserver/pkg/logf"
)

const (
	MAX_UPLOAD = 1000000000 // 1G
	//MAX_UPLOAD = 10000000 // 10mb

	FORM_NAME_FILE = "file"

	PATH_UPLOAD = "storage"

	//TYPES_FILES       = "/(gif|p?jpeg|(x-)?png)"
	TYPES_FILES_NOT = "/(%)"

	NOT_ACCEPT_FILE_TYPES = TYPES_FILES_NOT
)

var (
	imageTypes         = regexp.MustCompile(TYPES_FILES_NOT)
	notAcceptFileTypes = regexp.MustCompile(NOT_ACCEPT_FILE_TYPES)
)

type FileInfo struct {
	Key          string `json:"-"`
	ThumbnailKey string `json:"-"`
	Url          string `json:"url,omitempty"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Size         int64  `json:"size"`
	Error        string `json:"error,omitempty"`
	DeleteUrl    string `json:"deleteUrl,omitempty"`
	DeleteType   string `json:"deleteType,omitempty"`
}

//
//
//
func MethodUpload(w http.ResponseWriter, r *http.Request) {

	t1 := time.Now()

	//allow cross domain AJAX requests
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")

	ok, json, msg := UploadFiles(w, r)

	//
	//
	//
	JsonUp := []byte(json)

	//
	//
	//
	w.Header().Set("X-Custom-Header", "HeaderValue-"+cryptf.HashSeed(cryptf.Random(1, 30)))

	//
	//
	//
	w.Header().Set("Content-Type", "text/plain")

	if ok {

		w.WriteHeader(http.StatusOK)

	} else {

		w.WriteHeader(http.StatusBadRequest)
	}

	//
	//
	//
	w.Write(JsonUp)

	//
	//
	//
	t2 := time.Now()

	//
	//
	//
	logf.LogHandlerOn(w, r, msg, t1, t2)

}

//
// Allow only submissions containing multi / form-data, at this
// time it will not accept files coming in --data-binary.
// File size limitation must be defaulted on the platform
// multipart/form-data
// Do not accept
// application/x-www-form-urlencoded
// -H "Name-File: file1.jpg" \
// --data-binary "@files/file1.jpg"
// there are 2 possibilities for an upload
// one => send and replace
// two => keep 2
func UploadFiles(w http.ResponseWriter, r *http.Request) (bool, string, string) {

	var msgTmp string

	if strings.ToUpper(strings.TrimSpace(strings.Trim(r.Method, " "))) == "POST" {

		tmpVector := strings.SplitN(r.Header.Get("Content-Type"), ";", 2)
		contenType := tmpVector[0]

		contenType = strings.ToLower(strings.TrimSpace(strings.Trim(contenType, " ")))

		if contenType == "multipart/form-data" {

			ContentLength := r.ContentLength

			if ContentLength > MAX_UPLOAD {

				msgTmp = `Upload size exceeds maximum size of ` + fmt.Sprintf("%d", MAX_UPLOAD) + ` bytes`

				return false, `{"status":"error","msg":"` + msgTmp + `"}`, msgTmp
			}

			errup := r.ParseMultipartForm(MAX_UPLOAD)

			if errup != nil {

				msgTmp = fmt.Sprintf("%v", errup)
				return false, `{"status":"error","msg":"` + msgTmp + `"}`, msgTmp
			}

			//parsed multipart form
			multi := r.MultipartForm

			//fmt.Println("multi: ", multi.File)

			if len(multi.File[FORM_NAME_FILE+"[]"]) > 0 {

				// type of submission
				// 1 or 2, 1 => replaces the same file
				// 2 => create a new one from the same existing
				// 0 => or empty subistiui the same file
				//typeSend := r.FormValue("type_send")
				//fileUid := r.FormValue("file_uid")

				//fmt.Println("size array files: ", len(multi.File["files[]"]))
				//fmt.Printf("map %+v\n", multi.File)

				// Call our method to save the
				// sending of multiple files to disk
				//return UploadFormFileMulti(w, r, typeSend, fileUid, multi.File[FORM_NAME_FILE+"[]"])
				msgTmp = "multi uploads em desenvolviemnto"
				return true, `{"status":"ok","msg":"` + msgTmp + `"}`, msgTmp

			} else {

				// type of submission
				// 1 or 2, 1 => replaces the same file
				// 2 => create a new one from the same existing
				// 0 => or empty subistiui the same file
				typeSend := r.FormValue("type_send")
				fileUid := r.FormValue("file_uid")

				msgTmp = "uploads success"

				// Upload only 1 file,
				// saving a file to disk
				return UploadFormFileUnic(w, r, typeSend, fileUid)
				//return true, `{"status":"ok","msg":"` + msgTmp + `"}`, msgTmp
			}

		} else {

			/** @type {[type]} [description] */
			if contenType == "binary/octet-stream" {

				//
				//
				//
				//typeSend := r.Header.Get("type_send")

				//
				//
				//
				//fileUid := r.Header.Get("file_uid")

				msgTmp = "uploads Binario em desenvolvimento"

				//
				//
				//
				//return UploadFormFileBinary(w, r, typeSend, fileUid)
				//return true, "", ""
				return true, `{"status":"ok","msg":"` + msgTmp + `"}`, msgTmp

			} else {

				msgTmp = "Is not multipart/form-data"
				return false, `{"status":"ok","msg":"` + msgTmp + `"}`, msgTmp
			}

		}
	} else {

		msgTmp = "Is not POST"
		return false, `{"status":"ok","msg":"` + msgTmp + `"}`, msgTmp
	}
}

//
func (fi *FileInfo) ValidateNotType() (valid bool) {

	if notAcceptFileTypes.MatchString(fi.Type) {
		return true
	}

	fi.Error = "Filetype not allowed"
	return false
}

//
func (fi *FileInfo) ValidateSize() (valid bool) {

	if fi.Size > MAX_UPLOAD {

		fi.Error = "File is too big"

	} else {

		return true
	}

	return false
}

//
// testing if file exists
//
func ExistFileUpload(name string) bool {

	_, err := os.Stat(name)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

//
// deleting file from disk
//
func DeleteFileUpload(name string) bool {

	if ExistFileUpload(name) {

		err := os.Remove(name)

		if err != nil {

			return false

		} else {

			return true
		}
	} else {

		return false
	}
}
