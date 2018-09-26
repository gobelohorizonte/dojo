/*
* Go Library (C) 2017 Inc.
*
* @project     Ukkbox
* @package     main
* @author      @jeffotoni
* @size        01/06/2017
*
*
 */

package logf

import (
	//"github.com/jeffotoni/fileserver/models"
	"github.com/jeffotoni/fileserver/pkg/authi"
	"github.com/jeffotoni/fileserver/pkg/gcolor"
)

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const LayoutDateLog = "2006-01-02 15:04:05"
const LayoutDate = "2006-01-02"
const LayoutHour = "15:04:05"

//
//
//
func LogHandlerLoginOn(w http.ResponseWriter, r *http.Request, msg string, t1 time.Time, t2 time.Time, token string) {

	logHandler(w, r, "on", msg, t1, t2, token)
}

//
//
//
func LogHandlerOff(w http.ResponseWriter, r *http.Request, msg string, t1 time.Time, t2 time.Time) {

	logHandler(w, r, "off", msg, t1, t2, "")
}

//
//
//
func LogHandlerOn(w http.ResponseWriter, r *http.Request, msg string, t1 time.Time, t2 time.Time) {

	logHandler(w, r, "on", msg, t1, t2, "")
}

//
//
//
func logHandler(w http.ResponseWriter, r *http.Request, typeLog, msg string, t1 time.Time, t2 time.Time, token string) {

	var UserGlobal, UidUser, Uidwks, Expires, tokenUser string

	nameLog := "FileServer"

	DateHours := time.Now().Format(LayoutDateLog)

	nameLogBd := "FileServer"

	stringLog := gcolor.CyanCor("[" + nameLog + "]")

	stringLog = stringLog + " " + gcolor.RedCor(" | log: ok | ")

	timeI := fmt.Sprintf("%s", t1.Format(LayoutDateLog))

	timeSpent := fmt.Sprintf("%v", t2.Sub(t1))

	statusHttp := fmt.Sprintf("%v", http.StatusOK)

	remoteAddrTmp := fmt.Sprintf("%s", r.RemoteAddr)

	// clean
	remoteVet := strings.SplitN(remoteAddrTmp, ":", 2)

	remoteAddr := remoteVet[0]

	methodUrl := fmt.Sprintf("%s", r.Method)

	urlString := r.URL.String()

	stringLog = stringLog + " " + timeI

	stringLog = stringLog + " | "

	stringLog = stringLog + gcolor.YellowCor(statusHttp)

	stringLog = stringLog + " | "

	stringLog = stringLog + gcolor.BlueCor(timeSpent)

	stringLog = stringLog + "    | " + remoteAddr + "    | "

	stringLog = stringLog + gcolor.InBlueCor(methodUrl)

	stringLog = stringLog + "   "

	stringLog = stringLog + gcolor.GreenCor(urlString)

	//
	// screen data
	//
	// fmt.Println(stringLog)

	if token != "" {

		UserGlobal, UidUser, Uidwks, Expires = authi.ClaimsData(token)

	} else {

		tokenUser = authi.GetSplitTokenJwt(w, r)

		if tokenUser == "" {

			tokenUser = "--"
			UserGlobal = "--"
			UidUser = "--"

		} else {

			UserGlobal, UidUser, Uidwks, Expires = authi.ClaimsData(tokenUser)

			if UserGlobal == "" {

				UserGlobal = "--"
				UidUser = "--"
				Expires = "--"
			}
		}
	}

	fmt.Sprintf("%s", Uidwks)
	// fmt.Println("div: ", UserGlobal, ", ", UidUser, ", ", Expires)
	log.Println(typeLog, tokenUser, UserGlobal, UidUser, Expires, DateHours, nameLogBd, timeSpent, statusHttp, remoteAddr, methodUrl, urlString, msg)
}
