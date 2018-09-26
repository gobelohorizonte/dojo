package authi

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jeffotoni/fileserver/certs"
	"github.com/jeffotoni/fileserver/models"
	"github.com/jeffotoni/fileserver/repo"
)

import (
	"fmt"
	"net/http"
	"strings"
)

func ClaimsData(token string) (string, string, string, string) {

	if token != "" {
		// start
		parsedToken, err := jwt.ParseWithClaims(token, &models.Claim{}, func(*jwt.Token) (interface{}, error) {

			return certs.PublicKey, nil

		})

		if err != nil || !parsedToken.Valid {

			return "", "", "", ""
		}

		claims, ok := parsedToken.Claims.(*models.Claim)

		if !ok {

			return "", "", "", ""
		}

		// review implementation, performance
		// validate user before, is there a user?
		//if !DynamodbValidUserAccess(claims.User) {
		if !repo.PgUserValid(claims.User) {

			return "", "", "", ""
		}
		//review implementation, performance

		return claims.User, claims.Uid, claims.Uidwks, fmt.Sprintf("%v", claims.ExpiresAt)

	} else {

		return "", "", "", ""
	}
}

// TokenGlobal = token
// UserGlobal = claims.User
// ExpiGlobal = fmt.Sprintf("%d", claims.ExpiresAt)
// fmt.Println("User: ", claims.User)
// func2(w, r)
func GetSplitTokenJwt(w http.ResponseWriter, r *http.Request) string {

	var Authorization string

	Authorization = r.Header.Get("Authorization")

	if Authorization == "" {

		Authorization = r.Header.Get("authorization")
	}

	// browsers
	if Authorization == "" {

		Authorization = r.Header.Get("Access-Control-Allow-Origin")
	}

	if Authorization != "" {

		auth := strings.SplitN(Authorization, " ", 2)

		if len(auth) != 2 || strings.TrimSpace(strings.ToLower(auth[0])) != "bearer" {

			return ""
		}

		token := strings.Trim(auth[1], " ")
		token = strings.TrimSpace(token)

		return token
	} else {

		return ""
	}
}
