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
* @size        21/08/2017
*
 */

package util

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

///
func GetWdLocal(level int) string {

	pathNow, _ := os.Getwd()

	pathV := strings.Split(pathNow, "/")

	pathNew := pathV[:len(pathV)-level]

	pathNewOrg := strings.Join(pathNew, "/")

	return pathNewOrg
}

//  Exist only file
func FileExist(name string) bool {

	//if _, err := os.Stat(name); os.IsNotExist(err) {
	if stat, err := os.Stat(name); err == nil && !stat.IsDir() {

		return true
	}

	return false
}

func DirExist(path string) bool {

	//if _, err := os.Stat(path); err == nil {
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {

		return true
	}

	return false

	// if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
	// 	// path/to/whatever does not exist
	// }
}

func CreateDirIfNotExist(dir string) bool {

	if _, err := os.Stat(dir); os.IsNotExist(err) {

		err = os.MkdirAll(dir, 0755)

		if err != nil {
			//log.Save.Println(err)
			return false
		}
	}

	return true
}

//
func GetExistFile(pathNewOrg string) string {

	_, err := os.Stat(pathNewOrg)

	if err == nil {

		return pathNewOrg

	} else if os.IsNotExist(err) {

		return ""

	} else {

		return ""
	}
}

func UkkGwd(NameDir string, level int) string {

	pathNewOrg := GetWdLocal(level)

	pathNewOrg = pathNewOrg + "/" + NameDir

	return GetExistFile(pathNewOrg)
}

//
//
//
func ValidJson(w http.ResponseWriter, r *http.Request, bodyJson []byte) (bool, string) {

	if len(bodyJson) == 0 {

		return true, `{"status":"error","msg":"Missing Json"}`

	}

	//
	// Looking for keys in the first and last position
	//
	last_pos := len(bodyJson) - 1

	//
	//
	//
	if string(bodyJson[0]) != "{" {

		return true, `{"status":"error","msg":"Missing keys on your json '{'"}`

	} else if string(bodyJson[last_pos]) != "}" {

		return true, `{"status":"error","msg":"Missing keys on your json '}'"}`

	} else {

		return false, ""
	}
}

//
//
//
func ValidHeader(r *http.Request) string {

	//
	// Accept
	//
	contentType := r.Header.Get("Content-Type")

	////
	//fmt.Println(contentType)

	tmpContent := strings.ToLower(strings.TrimSpace(contentType))
	tmpContentV := strings.Split(tmpContent, ";")

	if tmpContentV[0] == "application/x-www-form-urlencoded" {

		return "FORM"

	} else if strings.ToLower(strings.TrimSpace(contentType)) == "application/x-www-form-urlencoded" {

		return "FORM"

	} else if strings.ToLower(strings.TrimSpace(contentType)) == "application/json" {

		return "JSON"

	} else {

		return "Error"
	}
}

// parseTimeStrict parses a formatted string and returns the time value it
// represents. The output is identical to time.Parse except it returns an
// error for strings that don't format to the input value.
//
// An example where the output differs from time.Parse would be:
// parseTimeStrict("1/2/06", "11/31/15")
//	- time.Parse returns "2015-12-01 00:00:00 +0000 UTC"
//	- parseTimeStrict returns an error
func DateValid(layout, value string) bool {

	t, err := time.Parse(layout, value)

	if err != nil {

		//println(err)
		return false
	}

	if t.Format(layout) != value {

		return false
	}

	return true
}

func RemoteAddr(r *http.Request) string {

	remoteAddrTmp := fmt.Sprintf("%s", r.RemoteAddr)
	// clean
	remoteVet := strings.SplitN(remoteAddrTmp, ":", 2)
	remoteAddr := remoteVet[0]

	return remoteAddr
}

func UserAgent(r *http.Request) string {

	return r.UserAgent()
}

func CheckErr(err error) {

	if err != nil {

		fmt.Println(err)
	}
}

//
// ex:
// ano := fmt.Sprintf("%s", fmt.Sprintf("%d", t.Year()))
// mes := fmt.Sprintf("%s", fmt.Sprintf("%d", Mes(fmt.Sprintf("%s", t.Month()))))
// dia := fmt.Sprintf("%s", fmt.Sprintf("%d", t.Day()))
func Mes(mes string) int {

	if len(mes) > 0 {
		var months = make(map[string]int)

		months["January"] = 1
		months["February"] = 2
		months["March"] = 3
		months["April"] = 4
		months["May"] = 5
		months["June"] = 6
		months["July"] = 7
		months["August"] = 8
		months["September"] = 9
		months["October"] = 10
		months["November"] = 11
		months["December"] = 12

		return months[mes]

	} else {

		return 0
	}
}

//
//
//
func ReverseString(input []string) []string {

	if len(input) == 0 {
		return input
	}

	return append(ReverseString(input[1:]), input[0])
}

//
//
//
func ReverseSort(input []string) []string {

	sort.Sort(sort.Reverse(sort.StringSlice(input)))
	return input
}
