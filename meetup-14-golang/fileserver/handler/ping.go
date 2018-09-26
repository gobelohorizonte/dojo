/*
* Go Library (C) 2017 Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
* @project     FileServer
* @package     main
* @author      @jeffotoni
* @size        23/09/2018
*
 */

package handler

import (
	"net/http"
	"time"
)

import (
	"github.com/jeffotoni/fileserver/pkg/cryptf"
	"github.com/jeffotoni/fileserver/pkg/logf"
)

//
// Testing whether the service is online
//
func MethodPing(w http.ResponseWriter, r *http.Request) {

	//
	//
	//
	t1 := time.Now()

	//
	//
	//
	json := `{"msg":"pong"}`

	//
	//
	//
	pong := []byte(json)

	//
	//
	//
	// w.Header().Set(HttpHeaderTitle, HttpHeaderMsg)

	//
	//
	//
	w.Header().Set("X-Custom-Header", "HeaderValue-"+cryptf.HashSeed(cryptf.Random(1, 30)))

	//
	//
	//
	w.Header().Set("Content-Type", "application/json")

	//
	//
	//
	w.WriteHeader(http.StatusOK)

	//
	//
	//
	w.Write(pong)

	//
	//
	//
	t2 := time.Now()

	//
	//
	//
	msg := "Ping success"

	//
	//
	//
	logf.LogHandlerOff(w, r, msg, t1, t2)

}
