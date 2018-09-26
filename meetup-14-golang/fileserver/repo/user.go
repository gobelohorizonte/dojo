package repo

import (
	"log"
	"strings"
)

import (
	"database/sql"
	"github.com/jeffotoni/fileserver/pkg/cryptf"
	"github.com/jeffotoni/fileserver/pkg/pg"
)

func PgUserValid(Email string) bool {

	// local
	var Db *sql.DB

	// Db...
	if interf := pg.PgConnect(pg.DB_NAME_1, pg.DB_HOST_1, pg.DB_USER_1, pg.DB_PASSWORD_1, pg.DB_PORT_1, pg.DB_SORCE_1); interf != nil {

		Db = interf.(*sql.DB)

	} else {

		log.Println("error ao fazer connect goPgUpdate com Db..", interf)
		return false
	}

	var logi_uid string

	// status
	// 1 = ative
	// 2 = inative
	logi_status := 1

	Email = strings.ToLower(strings.TrimSpace(Email))

	err := Db.QueryRow("SELECT logi_uid FROM fileserver_login WHERE lower(logi_email)=$1 AND logi_status=$2", Email, logi_status).Scan(&logi_uid)

	switch {

	case err == sql.ErrNoRows:

		return false

	case err != nil:

		return false

	default:

		return true
	}
}

// user authentication, the method returns user uid
// and its workspace, there are no users
// without workspaces
func PgAuthUser(Email, Password string) (string, string) {

	// local
	var Db *sql.DB

	// Db...
	if interf := pg.PgConnect(pg.DB_NAME_1, pg.DB_HOST_1, pg.DB_USER_1, pg.DB_PASSWORD_1, pg.DB_PORT_1, pg.DB_SORCE_1); interf != nil {

		Db = interf.(*sql.DB)

	} else {

		log.Println("error ao fazer connect goPgUpdate com Db..", interf)
		return "", ""
	}

	var logi_uid, logi_password, wks_uid string

	Email = strings.ToLower(strings.TrimSpace(Email))

	err := Db.QueryRow("SELECT logi_uid,logi_password,wks_uid FROM fileserver_login LEFT JOIN fileserver_workspace ON (wks_uiduser = logi_uid) WHERE lower(logi_email)=$1", Email).Scan(&logi_uid, &logi_password, &wks_uid)

	switch {

	case err == sql.ErrNoRows:

		return "", ""

	case err != nil:

		return "", ""

	default:

		if cryptf.CheckBlowfish(Password, logi_password) {

			return logi_uid, wks_uid

		} else {

			return "", ""
		}
	}
}
