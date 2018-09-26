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

/////////
package route

import (
	"net/http"
)

var (
	confServer *http.Server

	CorsAllow = []string{"http://localhost", "http://localhost/example", "http://localhost:3000"}

	CorsAllowedMethods = []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"}

	CorsAllowedHeaders = []string{"*"}
)

var (
	PORT_SERVER = "5000"
)

const (
	maxClients = 50000 // simultaneos

	NewLimiter = 20000 // 20k requests per second

	_      = iota
	KB int = 1 << (10 * iota)
	MB int = 100 << (10 * iota)
	GB int = 1000 << (10 * iota)

	Protocol = "http://"

	Schema = "localhost"

	// Schema = "34.225.180.116"

	HttpHeaderMsg = `Good Server, obrigado!`

	MaxHeaderByte = GB

	// pipng
	HandlerPing = "/v1/test/ping"

	// POST
	HandlerHello = "/v1/test/hello"

	// The login of the api, to work needs the access keys key_access and token_access,
	// that is the user registered in the platform will have the two keys available
	// to use the apis of the platform
	// POST
	HandlerLogin = "/v1/user/login"

	// Will be responsible for creating new account on the platform,
	// to use this handler requires the KEY_ACESS AND TOKEN_ACESS
	// You can only use the api if you have these 2 fields, only
	// these 2 fields will have the users that register in the SITE.
	// You will receive a post with the fields to create the user
	// and validate login / email and their respective fields.
	// HandlerCreateUser = "/create/user"
	// POST
	HandlerCreateUser = "/v1/user"

	//
	// The recover password will receive the account and
	// the respective keys of the site and send to the registered email
	//HandlerResetPassword = "/reset/password"
	// POST
	HandlerResetPassword = "/v1/user/reset_password"

	// It is responsible for allowing the sending of files to the server, to use these handler the user will have to have the token sent at the moment of login.
	// The header is used "Authorization: Bearer <token>" to allow sending.
	// The upload will be accepted multipart / form-data for single or multiple sending with multiple option.
	// Example in cURL would be -F | --form
	// The upload will also not accept the possibility of loading binary files at the first moment,
	// ie the cURL --data-binary option will not be allowed, due to the control limitations imposed by the protocol.
	// Example in cURL would be --data-binary
	//HandlerUpload = "/upload"
	// POST
	HandlerUpload = "/v1/file/upload"

	//
	// handler remove file, was set logically in the trash column warning the system
	// that it can no longer display the file on the screen, it went to the trash
	// HandlerUploadRemoveDefinitive = "/file/remove/definitive"
	// DELETE
	HandlerUploadRemoveDefinitive = "/v1/file"

	// this handler warns the system to remove from the trash
	// and physically from the system
	// POST
	HandlerUploadRemoveFileTrash = "/v1/file/remove_trash"

	// the user of the token itself can
	// authenticate and the validated
	// passer parameter is the key
	// received in the creation of the user
	// the method logically closes within
	// 24 hours the account will actually
	// be closed in every system
	// POST
	HandlerCloseAccount = "/v1/user/close_account"

	// Account restoration is given by e-mail,
	// he may restate the account up
	// to 24 hours before closing.
	// After this the account is finally closed
	// POST
	HandlerRestoreAccount = "/v1/user/restore_account"

	// only users of type admin or level 5
	// can enable and disable users,
	// the tokem is validated
	// to check the permission
	// POST
	HandlerDisableUser = "/v1/user/disable"

	// only users of type admin or level 5
	// can enable and disable users,
	// the tokem is validated
	// to check the permission
	// POST
	HandlerEnableUser = "/v1/user/enable"

	// the download will be done
	// from the physical disk or
	// the cloud storage
	// managed by java
	// GET
	HandlerDownload = "/v1/file/download"

	//
	//
	//
	// POST
	HandlerConfirmEmail = "/v1/user/confirm_email"
)

var (
	//
	//
	//
	HttpConfirmEmail = Protocol + Schema + ":" + PORT_SERVER + "" + HandlerConfirmEmail

	//
	//
	//
	HttpRestoreAccount = Protocol + Schema + ":" + PORT_SERVER + "" + HandlerRestoreAccount
)
