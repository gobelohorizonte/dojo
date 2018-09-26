/*
* FileServer Go Library (C) 2017 Inc.
*
* @project     FileServer
* @package     main
* @author      @jeffotoni
* @size        23/09/2018

* @description Our main auth will be responsible for validating our
* handlers, validating users and will also be in charge of creating users,
* removing them and doing their validation of access.
* We are using jwt to generate the tokens and validate our handlers.
* The logins and passwords will be in the AWS Dynamond database.
*
* $ openssl genrsa -out private.rsa 1024
* $ openssl rsa -in private.rsa -pubout > public.rsa.pub
*
 */

package auth

import (
	//"crypto/rsa"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/jeffotoni/fileserver/certs"
	"github.com/jeffotoni/fileserver/models"
	"github.com/jeffotoni/fileserver/pkg/authi"
	"github.com/jeffotoni/fileserver/pkg/logf"
	"github.com/jeffotoni/fileserver/repo"
	"github.com/jeffotoni/fileserver/util"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// expiracao em segundos
const secondExpires = 600

var (

	//pathPrivate = "../certs/private.rsa"
	//pathPublic  = "../certs/public.rsa.pub"

	ProjectTitle = "jwt FileServer"

	//
	HEDER_X_KEY = "fileserver2018golangbh"

	// Base64 Encode KEY
	HEDER_X_KEY_B64 = "ZmlsZXNlcnZlcjIwMThnb2xhbmdiaA=="

	// md5(basicHash) BD
	HEDER_X_KEY_RESTORE_ACCOUNT = "7199600eca3d5c4e4880c7950295dc41"

	// Base64 KEY
	HEDER_X_KEY_RESTORE_ACCOUNT_B64 = "ZmlsZXNlcnZlcjIwMThnb2xhbmdiaA=="
)

// discontinued, we will not use
// this more expensive implementation.
func GetDirKeys() (string, string) {

	var pathApp string

	pwd, err := os.Getwd()

	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	VetDir := strings.Split(pwd, "/")

	lenght := len(VetDir) - 2

	for i, path := range VetDir {

		if lenght != i {

			pathApp += fmt.Sprintf("%v%s", path, "/")

		} else {

			break
		}
	}

	pathApp = pathApp + "" + "certs/private.rsa"
	pathApp2 := pathApp + "" + "certs/public.rsa.pub"

	//fmt.Println(pathApp)

	return pathApp, pathApp2
}

//
// Structure of our server configurations
//
type JsonMsg struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
}

//
// jwt init
//
func init() {

	var errx error

	privateByte := []byte(certs.RSA_PRIVATE)

	certs.PrivateKey, errx = jwt.ParseRSAPrivateKeyFromPEM(privateByte)

	if errx != nil {

		// WriteJson("error", "Could not parse privatekey!")
		return
	}

	publicByte := []byte(certs.RSA_PUBLIC)

	//
	//
	//
	certs.PublicKey, errx = jwt.ParseRSAPublicKeyFromPEM(publicByte)

	if errx != nil {

		// WriteJson("error", "ould not parse publickey!")
		return
	}
}

//
// jwt 'GenerateJWT'
//
func GenerateJWT(model models.User) (string, string) {

	//
	// Generating date validation to return to the user
	//
	Expires := time.Now().Add(time.Second * secondExpires).Unix()

	//
	// convert int64
	//
	ExpiresInt64, _ := strconv.ParseInt(fmt.Sprintf("%v", Expires), 10, 64)

	//
	// convert time unix to Date RFC
	//
	ExpiresDateAll := time.Unix(ExpiresInt64, 0)

	// Data e Hora de expircao
	ExpiresDate := ExpiresDateAll.Format("2006-01-02 15:04:05")

	//
	// claims Token data, the header
	//
	claims := models.Claim{

		User: model.Login,

		Uid: model.Uid,

		Uidwks: model.Uidwks,

		StandardClaims: jwt.StandardClaims{

			//
			// Expires in 24 hours * 10 days
			//
			ExpiresAt: Expires,
			Issuer:    ProjectTitle,
		},
	}

	//
	// Generating token
	//
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	//
	// Transforming into string
	//
	tokenString, err := token.SignedString(certs.PrivateKey)

	if err != nil {

		return "Could not sign the token!", "2006-01-02"
	}

	//
	// return token string
	//
	return tokenString, ExpiresDate
}

//
// base64 (md5(key))
//
// login e password default in base 64
// curl -X POST -H "Content-Type: application/json"
// -H "Authorization: Basic ZTg5NjFlZDczYTQzMzE0YWYyY2NlNDdhNGY1YjY1ZGI=:ZGExMjRhMDAwNTE1MDUyYzFlNWJjNmU0NzQ4Yzc3ZTU="
// "https://localhost:9001/token"
//
func AuthBasicJwt(w http.ResponseWriter, r *http.Request) bool {

	//
	// do it now
	//
	if AuthBasicKey(w, r) {

		return true

	} else {

		return false
	}
}

//
//
//
func TokenClaimsJwt(w http.ResponseWriter, r *http.Request) string {

	return (authi.GetSplitTokenJwt(w, r))
}

//
//
//
func TokenJwtClaimsValid(w http.ResponseWriter, r *http.Request) bool {

	token := authi.GetSplitTokenJwt(w, r)

	if token != "" {

		// fmt.Println(token)
		// start
		parsedToken, err := jwt.ParseWithClaims(token, &models.Claim{}, func(*jwt.Token) (interface{}, error) {

			return certs.PublicKey, nil

		})

		if err != nil || !parsedToken.Valid {

			// HttpWriteJson(w, r, "error", "Your token has expired!", http.StatusAccepted)
			return false
		}

		_, ok := parsedToken.Claims.(*models.Claim)

		// fmt.Println("ok", ok)
		return ok

	} else {

		return false
	}
}

//
// ValidateHandler
//
func ValidateHandler(w http.ResponseWriter, r *http.Request) bool {

	if !TokenJwtClaimsValid(w, r) {

		//w.WriteHeader(http.StatusAccepted)
		// fmt.Fprintln(w, r, "There's something strange about your token!")
		//HttpWriteJson(w, r, "error", "There's something strange about your token!", http.StatusAccepted)
		return false
	}

	return true
}

//
// Returns json without typing in http
//
func WriteJson(Status string, Msg string) {

	msgJsonStruct := &JsonMsg{Status, Msg}

	msgJson, errj := json.Marshal(msgJsonStruct)

	if errj != nil {

		fmt.Println(`{"status":"error","msg":"We could not generate the json error!"}`)
		return
	}

	fmt.Println(msgJson)
}

//
//
//
func GetXKey(w http.ResponseWriter, r *http.Request) string {

	//
	//
	// 'X-Key: YOUR-API-KEY-HERE'
	//
	auth := strings.SplitN(r.Header.Get("X-Key"), " ", 2)

	//
	//
	//
	if len(auth) <= 0 {

		//HttpWriteJson(w, r, "error", "Your X-Key it's wrong!", http.StatusAccepted)
		return ""
	}

	//
	//
	//
	tokenBase64 := strings.Trim(auth[0], " ")

	//
	//
	//
	tokenBase64 = strings.TrimSpace(tokenBase64)

	//
	// User, Login byte
	//
	tokenUserDecode, _ := b64.StdEncoding.DecodeString(tokenBase64)

	//
	// User, Login string
	//
	tokenUserDecodeS := strings.TrimSpace(strings.Trim(string(tokenUserDecode), " "))

	return tokenUserDecodeS
}

//
//
//
func GetXKeyUrl(w http.ResponseWriter, r *http.Request) string {

	//
	//
	// 'X-Key: YOUR-API-KEY-HERE'
	//
	auth := r.FormValue("k")

	//
	//
	//
	if auth == "" {

		//HttpWriteJson(w, r, "error", "Your X-Key it's wrong!", http.StatusAccepted)
		return ""
	}

	//
	//
	//
	tokenBase64 := strings.Trim(auth, " ")

	//
	//
	//
	tokenBase64 = strings.TrimSpace(tokenBase64)

	//
	// User, Login byte
	//
	tokenUserDecode, _ := b64.StdEncoding.DecodeString(tokenBase64)

	//
	// User, Login string
	//
	tokenUserDecodeS := strings.TrimSpace(strings.Trim(string(tokenUserDecode), " "))

	return tokenUserDecodeS
}

//
// base64 (md5(key))
//
// login e password default in base 64
// curl -X POST -H "Content-Type: application/json"
// -H "Authorization: Basic ZTg5NjFlZDczYTQzMzE0YWYyY2NlNDdhNGY1YjY1ZGI=:ZGExMjRhMDAwNTE1MDUyYzFlNWJjNmU0NzQ4Yzc3ZTU="
// "https://localhost:9001/token"
//
func AuthBasicKey(w http.ResponseWriter, r *http.Request) bool {

	tokenUserDecodeS := GetXKey(w, r)

	fmt.Println(tokenUserDecodeS + " :: " + HEDER_X_KEY)
	//
	// Validate user and password in the database
	//
	if tokenUserDecodeS == HEDER_X_KEY {

		return true

	} else {

		//HttpWriteJson(w, r, "error", "Your X-Key it's wrong ...", http.StatusAccepted)
		return false
	}

	//defer r.Body.Close()
	//HttpWriteJson(w, r, "error", "Your X-Key it's wrong ..", http.StatusAccepted)
	return false
}

//
// base64 (md5(key))
//
// login e password default in base 64
// curl -X POST -H "Content-Type: application/json"
// -H "Authorization: Basic ZTg5NjFlZDczYTQzMzE0YWYyY2NlNDdhNGY1YjY1ZGI=:ZGExMjRhMDAwNTE1MDUyYzFlNWJjNmU0NzQ4Yzc3ZTU="
// "https://localhost:9001/token"
//
func GetAuthBasicKey(w http.ResponseWriter, r *http.Request) bool {

	tokenUserDecodeS := GetXKeyUrl(w, r)

	//
	// Validate user and password in the database
	//
	if tokenUserDecodeS == HEDER_X_KEY {

		return true

	} else {

		//HttpWriteJson(w, r, "error", "Your X-Key it's wrong ...", http.StatusAccepted)
		return false
	}

	//defer r.Body.Close()
	//HttpWriteJson(w, r, "error", "Your X-Key it's wrong ..", http.StatusAccepted)
	return false
}

//
// base64 (md5(key))
//
// login e password default in base 64
// curl -X POST -H "Content-Type: application/json"
// -H "Authorization: Basic ZTg5NjFlZDczYTQzMzE0YWYyY2NlNDdhNGY1YjY1ZGI=:ZGExMjRhMDAwNTE1MDUyYzFlNWJjNmU0NzQ4Yzc3ZTU="
// "https://localhost:9001/token"
//
func GetAuthBasicKeyRestoreAccount(w http.ResponseWriter, r *http.Request) bool {

	tokenUserDecodeS := GetXKeyUrl(w, r)

	//
	// Validate user and password in the database
	//
	if tokenUserDecodeS == HEDER_X_KEY_RESTORE_ACCOUNT {

		return true

	} else {

		//HttpWriteJson(w, r, "error", "Your X-Key it's wrong ...", http.StatusAccepted)
		return false
	}

	//defer r.Body.Close()
	//HttpWriteJson(w, r, "error", "Your X-Key it's wrong ..", http.StatusAccepted)
	return false
}

//
// Returns json by typing on http
//
func HttpWriteJsonNew(w http.ResponseWriter, r *http.Request, Status string, Msg string, httpStatus int) {

	//
	//
	//
	t1 := time.Now()

	//
	//
	//
	msg := Msg

	//
	//
	//
	msgJsonStruct := &JsonMsg{Status, Msg}

	//
	//
	//
	msgJson, errj := json.Marshal(msgJsonStruct)

	//
	//
	//
	if errj != nil {

		fmt.Fprintln(w, `{"status":"error","msg":"We could not generate the json error!"}`)
		return
	}

	//
	//
	//
	w.WriteHeader(httpStatus)

	//
	//
	//
	w.Header().Set("Content-Type", "application/json")

	//
	//
	//
	w.Write(msgJson)

	//
	//
	//
	t2 := time.Now()

	//
	//
	//
	logf.LogHandlerOff(w, r, msg, t1, t2)
	//log.Println(msg, t1, t2)
}

////////// start
///

//
// Testing whether the service is online
//
func MethodLogin(w http.ResponseWriter, r *http.Request) {

	TokenLocal := ""

	t1 := time.Now()

	statusOk := false

	ok, emailUser, uidUser, uidWks := ValidUser(w, r)

	if ok {

		// do it now
		var model models.User

		model.Login = emailUser

		model.Uid = uidUser

		model.Uidwks = uidWks

		model.Password = ""

		model.Role = "user-default"

		token, expires := GenerateJWT(model)

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
	msg := "login Error authentication"

	//
	//
	//
	logf.LogHandlerOff(w, r, msg, t1, t2)
}

///////// user
///

type UserValid struct {
	Email string `json:"user"`
	//Password string `json:"password"`
	Password json.Number `json:"password,Number"`
}

//
// Create User Form
//
func ValidUser(w http.ResponseWriter, r *http.Request) (bool, string, string, string) {

	//
	//
	//
	var Email, Password2 string

	Password := UserValid{}

	// just scape
	fmt.Sprint("%v", Password)

	//
	//
	//
	typeHeader := util.ValidHeader(r)

	//
	//
	//
	if typeHeader == "JSON" {

		//
		// Validating json if correct
		//
		bodyJson, _ := ioutil.ReadAll(r.Body)

		//
		//
		//
		jsonValid, msgJson := util.ValidJson(w, r, bodyJson)

		//
		// true error
		// false success
		//
		if jsonValid {

			return false, msgJson, "", ""

		} else {

			//
			//
			//
			var L UserValid

			//
			//
			//
			err := json.Unmarshal(bodyJson, &L)

			if err != nil {

				return false, `{"status":"error","msg":"Json did not work: ` + fmt.Sprintf("%s", err) + `"}`, "", ""
			}

			//
			//
			//
			Email = L.Email
			Password := L.Password
			Password2 = fmt.Sprintf("%s", Password)
		}

	} else if typeHeader == "FORM" {

		// Email
		// Password
		Email = r.FormValue("user")
		Password := r.FormValue("password")
		Password2 = fmt.Sprintf("%s", Password)

	} else {

		//error
		return false, `{"status":"error","msg":"Invalid Content-Type"}`, "", ""
	}

	// Password2 := fmt.Sprintf("%s", Password)
	//
	//
	//
	Password2 = strings.TrimSpace(strings.Trim(Password2, " "))

	//
	// emails valid
	//
	Email = strings.ToLower(strings.TrimSpace(strings.Trim(Email, " ")))

	if Email == "" {

		return false, `{"status":"error","msg":"Empty field [User]"}`, "", ""

	} else if Password2 == "" {

		return false, `{"status":"error","msg":"Empty field [Password]"}`, "", ""
	}

	if len(Email) >= 201 {

		return false, `{"status":"error","msg":"Very large size, allowed 200 characters"}`, "", ""

	} else if len(Password2) >= 101 {

		return false, `{"status":"error","msg":"Very large size, allowed 100 characters"}`, "", ""
	}

	//
	// Is the user active?
	//
	//if !DynamodbValidUserAccess(Email) {
	if !repo.PgUserValid(Email) {

		return false, `{"status":"error","msg":"Wrong user is not active!"}`, "", ""
	}

	//
	// Validate if the user exists
	//
	uidUser, uidWrks := repo.PgAuthUser(Email, Password2)

	if uidUser == "" {

		return false, `{"status":"error","msg":"Wrong password, try again!"}`, "", ""
	}

	return true, Email, uidUser, uidWrks
}
