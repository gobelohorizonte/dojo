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

/*
* Go Library (C) 2017 Inc.
*
* @project     FileServer
* @package     main
* @author      @jeffotoni
* @size        23/09/2018
*
* @description Our route is responsible for making routes
* and defining authentication in the handlers and determining
* access permissions, we use jwt to generate the tokens and
* authorize the use of api.
*
* $ openssl genrsa -out private.rsa 1024
* $ openssl rsa -in private.rsa -pubout > public.rsa.pub
*
 */

///////// start api
package route

import (
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"

	"github.com/jeffotoni/fileserver/handler"
	auth "github.com/jeffotoni/fileserver/pkg/auth"
	cors "github.com/jeffotoni/fileserver/pkg/cors"
)

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	//"strings"
	"time"
)

//
// Type responsible for defining a function that returns boolean
//
type fn func(w http.ResponseWriter, r *http.Request) bool

type Adapter func(http.Handler) http.Handler

// Adapt h with all specified adapters.
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

//
// Function responsible for abstraction and receive the
// authentication function and the handler that will execute if it is true
//
func HandlerFuncAuth(authv fn, handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if authv(w, r) {

			handler(w, r)

		} else {

			auth.HandlerError(w, r)
		}
	}
}

func maxClientsFunc(h http.Handler, n int) http.HandlerFunc {

	sema := make(chan struct{}, n)

	return func(w http.ResponseWriter, r *http.Request) {

		sema <- struct{}{}

		defer func() { <-sema }()

		h.ServeHTTP(w, r)
	}
}

func Handlers() {

	// &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour}
	limiter := tollbooth.NewLimiter(NewLimiter, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})

	limiter.SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}).
		SetMethods([]string{"GET", "POST"})
	// Limit only GET and POST requests.
	//limiter.Methods = []string{"GET", "POST"}

	//
	//
	//
	mux := http.NewServeMux()

	//
	// defining Cors that can access our system
	// even though they are not being managed by the api or outside the domain
	corsx := cors.New(cors.Options{

		AllowedOrigins:   CorsAllow,
		AllowedMethods:   CorsAllowedMethods,
		AllowedHeaders:   CorsAllowedHeaders,
		AllowCredentials: true,
	})

	// cors allow
	cors.AllowAll().Handler(mux)

	// cors mux
	handlerCors := corsx.Handler(mux)

	// Ping testando API
	handlerPing := http.HandlerFunc(handler.MethodPing)

	//
	// Test api width ping
	// Public widthout X-Key
	//
	mux.Handle(HandlerPing, tollbooth.LimitFuncHandler(limiter, maxClientsFunc(handlerPing, maxClients)))

	//
	// Off the default mux
	// Does not need authentication, only user key and token
	// public / generate jwt
	//
	mux.Handle(HandlerLogin, tollbooth.LimitFuncHandler(limiter, HandlerFuncAuth(auth.AuthBasicJwt, handler.MethodLogin)))

	// //
	// // Private
	// // Test the api using the access keys
	// // ValidHandler
	mux.Handle(HandlerHello, tollbooth.LimitFuncHandler(limiter, HandlerFuncAuth(auth.ValidateHandler, handler.MethodHello)))

	//
	// Private
	// Upload files
	// ValidHandler
	mux.Handle(HandlerUpload, tollbooth.LimitFuncHandler(limiter, HandlerFuncAuth(auth.ValidateHandler, handler.MethodUpload)))

	// //
	// //
	// //
	// mux.Handle(HandlerDownload, tollbooth.LimitFuncHandler(limiter, HandlerFuncAuth(auth.ValidateHandler, MethodDownload)))

	//
	//
	//
	confServer = &http.Server{

		Addr: ":" + PORT_SERVER,

		Handler: handlerCors,
		//ReadTimeout:    30 * time.Second,
		//WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: MaxHeaderByte, // Size accepted by package

		ReadTimeout: 5 * time.Second,

		WriteTimeout: 10 * time.Second,
	}

	// start
	// service
	go func() {

		// service connections
		if err := confServer.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)

		}
	}()

	// Wait for interrupt signal
	// to gracefully shutdown
	// the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	// context timeout, cancel
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// shutdown CTRL+C
	if err := confServer.Shutdown(ctx); err != nil {

		log.Fatal("Server Shutdown:", err)
	}

	// log exist
	log.Println("Server exist")

	// listenAndServer
	log.Fatal(confServer.ListenAndServe())

}
