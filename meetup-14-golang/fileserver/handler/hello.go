/*
* FileServer Go Library (C) 2017 Inc.
*
* @project     FileServer
* @package     main
* @author      @jeffotoni
* @size        23/09/2018
 */

package handler

import (
	"github.com/jeffotoni/fileserver/pkg/cryptf"
	"github.com/jeffotoni/fileserver/pkg/logf"
)

import (
	"net/http"
	"time"
)

//
// Testing whether the service is online with X-key
func MethodHello(w http.ResponseWriter, r *http.Request) {

	//
	//
	//
	t1 := time.Now()

	//
	//
	//
	json := `{"status":"ok","msg":"Hello, welcome."}`

	//
	//
	//
	hello := []byte(json)

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
	w.Write(hello)

	//
	//
	//
	t2 := time.Now()

	//
	//
	//
	msg := "Hello success"

	//
	//
	//
	logf.LogHandlerOn(w, r, msg, t1, t2)

}
