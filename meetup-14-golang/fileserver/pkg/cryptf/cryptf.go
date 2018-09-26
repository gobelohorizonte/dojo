/*
* Go Library (C) 2017 Inc.
*
* @project     Ukkbox
* @package     main
* @author      @jeffotoni
* @size        01/06/2017
*
*
*
 */

package cryptf

//
//
//
import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"math/rand"
	"os"
	"time"
)

var (
	HASH_SALT      = "jeff2018&*$#@!yrtx.#568xejeszlx"
	SHA1_SALT      = "#@%@$@.w8w.xow8e.9e8x3**3x3@#x#"
	letters        = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	UKK_KEY_CIPHER = []byte("DKYPENJXW43SMTYCU6F5TMFVOUANMJNL")
)

//
//
//
func HashSeed(n int) string {

	b := make([]rune, n)
	for i := range b {

		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

//
//
//
func Random(min, max int) int { rand.Seed(time.Now().Unix()); return rand.Intn(max-min) + min }

// generating single seed, can not repeat, even calling in almost
// the same time interval
func RandUid() string {

	// generate 64 bits timestamp
	unix64bits := uint64(time.Now().UTC().UnixNano())

	buff := make([]byte, 128)

	numRead, err := rand.Read(buff)

	if numRead != len(buff) || err != nil {
		panic(err)
	}

	unixUid := fmt.Sprintf("%x", unix64bits)

	//fmt.Println(unixUid)
	//return unixUid

	unixUid = GSha1(Blowfish(unixUid))
	return unixUid
}

//
//
//
func Blowfish(password string) string {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {

		panic(err)
	}

	return string(bytes)
}

//
//
//
func CheckBlowfish(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//
//
//
func GHash256(password string) string {

	pass := password + HASH_SALT

	h := sha256.New()
	h.Write([]byte(pass))

	return fmt.Sprintf("%x", h.Sum(nil))
}

//
//
//
func GHash512(password string) string {

	pass := password + HASH_SALT

	h := sha512.Sum512([]byte(pass))

	return fmt.Sprintf("%x", h)
}

//
//
//
func GSha1(key string) string {

	data := []byte(key + SHA1_SALT)
	return (fmt.Sprintf("%x", sha1.Sum(data)))
}

//
func Md5(text string) string {

	h := md5.New()
	io.WriteString(h, text)
	return (fmt.Sprintf("%x", h.Sum(nil)))
}

//
func HashFile(filePath string) (string, error) {

	var returnMD5String string

	file, err := os.Open(filePath)
	if err != nil {

		return returnMD5String, err
	}

	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {

		return returnMD5String, err
	}

	hashInBytes := hash.Sum(nil)[:16]

	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}

//
func FSBase64Encode(textString string) string {

	text := []byte(textString)

	sEnc := b64.URLEncoding.EncodeToString(text)

	return sEnc
}

//
func FSBase64Decode(textString string) string {

	//text := []byte(textString)

	sDec, _ := b64.URLEncoding.DecodeString(textString)

	return string(sDec)
}
