/*
* Go Library (C) 2017 Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
* @project     Ukkbox
* @package     main
* @author      @jeffotoni
* @size        29/08/2017
*
 */

package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

import (
	"github.com/jeffotoni/fileserver/pkg/authi"
	"github.com/jeffotoni/fileserver/pkg/cryptf"
	"github.com/jeffotoni/fileserver/pkg/logf"
	"github.com/jeffotoni/fileserver/repo"
	"github.com/jeffotoni/fileserver/util"
)

// the multipart / form-data format, ie comes from a form
func UploadFormFileUnic(w http.ResponseWriter, r *http.Request, typeSend, fileUid string) (bool, string, string) {

	// start time
	t1 := time.Now()

	var Up repo.PgWorkUpload

	// typeSend and fileUid
	// type of submission:
	// 1, 2 or 0 ,
	// 1 => replaces the same file,
	// 2 => create a new one from the same existing
	// 0 => or empty subistiui the same file

	var msgTmp, nameFileUidUser, msgupload string
	var boolFileExist, dyupok bool

	// Looking for the file in the FormFile method
	file, handler, errf := r.FormFile(FORM_NAME_FILE)

	if errf != nil {

		msgTmp = fmt.Sprintf("%v", errf.Error())
		return false, `{"status":"error","msg":"` + msgTmp + `"}`, msgTmp

	}

	// close file
	defer file.Close()

	fi := &FileInfo{
		Name: handler.Filename,
		Type: handler.Header.Get("Content-Type"),
	}

	if !fi.ValidateSize() {

		msgTmp = "size not allowed"
		return false, `{"status":"error","msg":"` + msgTmp + `"}`, msgTmp
	}

	// validating type not allowed
	if fi.ValidateNotType() {

		msgTmp = "tye not allowed " + fi.Type
		return false, `{"status":"error","msg":"` + msgTmp + `"}`, msgTmp
	}

	upload_name := handler.Filename

	// fetching the user's token online
	tokenUser := authi.GetSplitTokenJwt(w, r)

	// seeking uid, and the user of claims
	UserEmail, UidUser, Uidwks, _ := authi.ClaimsData(tokenUser)

	//returns error if not found
	if UidUser == "" {

		msgTmp = "I was unable to retrieve the user's UID"
		return false, `{"status":"error","msg":"` + msgTmp + `"}`, msgTmp
	}

	// empty folder with white space means
	// it will stay in user dashbord
	upload_folder := repo.PgGetUploadUidFolder(UidUser)

	// trying to grab the extension
	//tmpv := strings.SplitN(handler.Filename, ".", 2)
	//upload_type := tmpv[1]
	upload_type := fi.Type

	// File Exists ?
	// uidFileExist := DyGetFileNameUpload(UidUser, UserEmail, upload_name)
	uidFileExist := repo.PgUploadNameExist(UidUser, upload_name, upload_folder)

	//nameFileUidUser = RandUid()

	// file exists
	// updates the existing one 1 or 0
	// create again if typeSend 2
	if uidFileExist != "" {

		if uidFileExist == "error" {

			msgTmp = `We did not find the user [` + UserEmail + `] in the database`
			return false, `{"status":"error","msg":"` + msgTmp + `"}`, msgTmp

		} else if typeSend == "1" || typeSend == "0" || typeSend == "" {

			// update file
			// receiving uid from the file
			nameFileUidUser = uidFileExist
			boolFileExist = true

		} else if typeSend == "2" { // create again if typeSend 2

			// we will implement in a second moment
			// create new
			// nameFileUidUser = RandUid()
			// boolFileExist = false

			nameFileUidUser = uidFileExist
			boolFileExist = true

		} else {

			// update file
			// receiving uid from the file
			nameFileUidUser = uidFileExist
			boolFileExist = true

		}
	} else {

		// uid, single for each upload, NEW FILE
		// will be the uid of the file registry
		nameFileUidUser = cryptf.RandUid()
		boolFileExist = false
	}

	fmt.Println("path local: ", os.Getwd()+" :::path-upload::: "+PATH_UPLOAD)

	// os.Exit(0)
	if !util.CreateDirIfNotExist(PATH_UPLOAD) {

		msgTmp = "nao conseguimos criar o dreitorio " + PATH_UPLOAD + " para armazenar os arquivos."
		return false, `{"status":"error","msg":"` + msgTmp + `"}`, msgTmp
	}

	///create dir to key
	pathUpKeyUser := PATH_UPLOAD + "/" + UidUser

	existPath, _ := os.Stat(pathUpKeyUser)

	// exist o path file
	// not exist create
	if existPath == nil {

		// create path
		os.MkdirAll(pathUpKeyUser, 0777)
	}

	// path physicist
	pathUserAcess := pathUpKeyUser + "/" + nameFileUidUser

	// create files
	f, erros := os.OpenFile(pathUserAcess, os.O_WRONLY|os.O_CREATE, 0777)

	// erros OpenFile
	if erros != nil {

		msgTmp = "We could not create the directory and file for your upload"
		return false, `{"status":"error","msg":"` + msgTmp + `"}`, msgTmp
	}

	// close file
	defer f.Close()

	// Copying the FormFile file
	// to our local disk file
	sizef, errcy := io.Copy(f, file)

	if errcy != nil {

		msgTmp = "We could not upload copy"
		return false, `{"status":"error","msg":"` + msgTmp + `"}`, msgTmp
	}

	upload_uid := nameFileUidUser
	upload_size := fmt.Sprintf("%d", sizef)

	//DateHours := time.Now().Format(LayoutDateLog)
	remoteAddrTmp := fmt.Sprintf("%s", r.RemoteAddr)

	// clean leaving only ip
	remoteVet := strings.SplitN(remoteAddrTmp, ":", 2)
	remoteAddr := remoteVet[0]

	DateUp := time.Now().Format(logf.LayoutDate)
	HoursUp := time.Now().Format(logf.LayoutHour)

	// // trying to capture the file's type
	// buffer := make([]byte, 512)
	// n, _ := file.Read(buffer)
	// // capturing file type
	// upload_type := fmt.Sprintf("%s", http.DetectContentType(buffer[:n]))

	//
	t2 := time.Now()

	timeSpent := fmt.Sprintf("%v", t2.Sub(t1))

	// Update file exist UID
	if boolFileExist {

		Up.Work_uid = upload_uid
		Up.Work_uiduser = UidUser
		Up.Work_name = upload_name
		Up.Work_uidwks = Uidwks
		Up.Work_user_up = UidUser
		Up.Work_user_criacao = UidUser
		Up.Work_ip = remoteAddr
		Up.Work_browser = r.UserAgent()
		Up.Work_size = upload_size
		Up.Work_data_up = DateUp
		Up.Work_hora_up = HoursUp
		Up.Work_timespent = timeSpent
		Up.Work_uidfolder = upload_folder

		// Update Register upload
		// dyupok, msgupload = DynamodbUpdateUpload(UidUser, upload_name, DateHours, timeSpent, upload_size)
		dyupok, msgupload = Up.PgUploadUpdate()

		if !dyupok {

			msgTmp = "error trying to update your database: [" + msgupload + "]"
			return false, msgupload, msgTmp
		} else {

			msgTmp = "update foi feito com sucesso!"
		}

	} else {

		Up.Work_uid = upload_uid
		Up.Work_uiduser = UidUser
		Up.Work_uidfolder = upload_folder
		Up.Work_name = upload_name
		Up.Work_uidwks = Uidwks
		Up.Work_user_up = UidUser
		Up.Work_user_criacao = UidUser
		Up.Work_ip = remoteAddr
		Up.Work_browser = r.UserAgent()
		Up.Work_type = upload_type
		Up.Work_size = upload_size
		Up.Work_timespent = timeSpent

		//fmt.Println(Uidwks)

		// Create Register upload
		// dyupok, _, msgupload = DynamodbUpload(upload_uid, tokenUser, UidUser, DateHours, timeSpent, remoteAddr, upload_folder, upload_name, upload_size, upload_type)
		dyupok, msgupload = Up.PgUploadInsert()

		// removing only if it is inserted
		if !dyupok {

			//remove file
			DeleteFileUpload(pathUserAcess)

			msgTmp = "error trying to update your database: [" + msgupload + "]"
			//return false, `{"status":"error","msg":"` + msgTmp + `"}`, msgTmp
			return false, msgupload, msgTmp

		} else {

			msgTmp = "novo upload foi feito com sucesso!"
		}
	}

	// fmt.Println("typeSend:", typeSend)
	// fmt.Println("UidExist ?:", uidFileExist)
	// fmt.Println("bool exist ?:", boolFileExist)

	// fmt.Println("UidUser:", UidUser)
	// fmt.Println("name file:", upload_name)
	// fmt.Println("Uid name:", nameFileUidUser)
	// fmt.Println("Path:", pathUserAcess)

	return true, `{"status":"ok","msg":"` + msgTmp + `","size":"` + fmt.Sprintf("%d", sizef) + `","name":"` + handler.Filename + `","id":"` + upload_uid + `"}`, msgTmp
}
