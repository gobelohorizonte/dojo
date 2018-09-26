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
	"github.com/jeffotoni/fileserver/models"
	"github.com/jeffotoni/fileserver/pkg/auth"
	"github.com/jeffotoni/fileserver/pkg/logf"
	"net/http"
	"time"
)

func MethodLogin(w http.ResponseWriter, r *http.Request) {

	TokenLocal := ""

	t1 := time.Now()

	statusOk := false

	ok, emailUser, uidUser, uidWks := auth.ValidUser(w, r)

	if ok {

		// do it now
		var model models.User

		model.Login = emailUser

		model.Uid = uidUser

		model.Uidwks = uidWks

		model.Password = ""

		model.Role = "user-default"

		token, expires := auth.GenerateJWT(model)

		if token == "" || expires == "" {

			jsonByte := []byte(`{"status":"error","msg":"Token error generating!"}`)

			w.WriteHeader(http.StatusOK)

			w.Header().Set("Content-Type", "application/json")

			w.Write(jsonByte)

		} else {

			//send email emailUser

			// write client
			TokenLocal = token
			statusOk = true

			jsonByte := []byte(`{"status":"ok","msg":"success","token":"` + token + `","expires":"` + expires + `"}`)

			w.WriteHeader(http.StatusOK)

			w.Header().Set("Content-Type", "application/json")

			w.Write(jsonByte)
		}

	} else {

		// Can contain a json as a message when
		// it gives some kind of error
		jsonByte := []byte(emailUser)

		w.WriteHeader(http.StatusOK)

		w.Header().Set("Content-Type", "application/json")

		w.Write(jsonByte)
	}

	//
	//
	//
	t2 := time.Now()

	if statusOk {

		msg := "login success"
		logf.LogHandlerLoginOn(w, r, msg, t1, t2, TokenLocal)

	} else {

		msg := "Error while trying to login"
		logf.LogHandlerOff(w, r, msg, t1, t2)
	}
}

//
// Testing whether the service is online
//
func HandlerError(w http.ResponseWriter, r *http.Request) {

	//
	//
	//
	t1 := time.Now()

	jsonByte := []byte(`{"status":"error","msg":"Error in your authentication!"}`)

	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonByte)

	//
	//
	//
	t2 := time.Now()

	//
	//
	//
	msg := "login Error authentication!"

	//
	//
	//
	logf.LogHandlerOff(w, r, msg, t1, t2)
}
