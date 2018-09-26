/*
* Go Library (C) 2017 Inc.
*
* @project     FileServer
* @package     main
* @author      @jeffotoni
* @size        23/09/2018
 */

package main

//
//
//
import (
	route "github.com/jeffotoni/fileserver/route"
	"os"
)

// inicializando
// com PORTA
// que desejar
func init() {

	////////// inicio
	port_tmp := os.Getenv("PORT_SERVER")

	if port_tmp != "" {

		route.PORT_SERVER = port_tmp

	} else {

		//if for argumentos OK
		if len(os.Args) == 2 && os.Args[1] != "" {

			route.PORT_SERVER = os.Args[1]

		}
	}
}

//
// start
//
func main() {

	//
	//
	//
	route.ShowScreenMain()

	//
	//
	//
	route.ShowScreen()

	//
	//
	//
	route.Handlers()
}
