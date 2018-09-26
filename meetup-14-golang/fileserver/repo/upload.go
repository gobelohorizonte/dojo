/*
* Go Library (C) 2017 Inc.
*
* @project     Ukkbox
* @package     main
* @author      @jeffotoni
* @size        16/07/2017
*
 */

package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jeffotoni/fileserver/pkg/pg"
	"log"
)

type PgWorkUpload struct {

	//Logi_id           string `json:"logi_id"`
	Work_id json.Number `json:"work_id,Number"`

	Work_uid             string `json:"work_uid"`
	Work_uiduser         string `json:"work_uiduser"`
	Work_uidfolder       string `json:"work_uidfolder"`
	Work_name            string `json:"work_name"`
	Work_uidwks          string `json:"work_uidwks"`
	Work_data_criacao    string `json:"work_data_criacao"`
	Work_data_up         string `json:"work_data_up"`
	Work_hora_criacao    string `json:"work_hora_criacao"`
	Work_hora_up         string `json:"work_hora_up"`
	Work_user_up         string `json:"work_user_up"`
	Work_user_criacao    string `json:"work_user_criacao"`
	Work_ip              string `json:"work_ip"`
	Work_browser         string `json:"work_browser"`
	Work_type            string `json:"work_type"`
	Work_size            string `json:"work_size"`
	Work_timespent       string `json:"work_timespent"`
	Work_descricao       string `json:"work_descricao"`
	Work_remocao         string `json:"work_remocao"`
	Work_remocao_data    string `json:"work_remocao_data"`
	Work_remocao_hora    string `json:"work_remocao_hora"`
	Work_remocao_user    string `json:"work_remocao_user"`
	Work_remocao_ip      string `json:"work_remocao_ip"`
	Work_remocao_browser string `json:"work_remocao_browser"`
}

// upload_uid, tokenUser, uidLogin, DateHours, timeSpent, remoteAddr,
// upload_folder, upload_name, upload_size, upload_type
// start struct user
// var Upload PgWorkUpload
// fi_xxxxxxxxxxxxxxxxxxxuid
func (Up *PgWorkUpload) PgUploadInsert() (bool, string) {

	// start connect
	// Postgresql
	// local
	var Db *sql.DB

	// Db...
	if interf := pg.PgConnect(pg.DB_NAME_1, pg.DB_HOST_1, pg.DB_USER_1, pg.DB_PASSWORD_1, pg.DB_PORT_1, pg.DB_SORCE_1); interf != nil {

		Db = interf.(*sql.DB)

	} else {

		log.Println("error ao fazer connect PgUploadInsert com Db..", interf)
		return false, "error ao fazer connect PgUploadInsert"
	}

	nameTableUpload := "fi_" + Up.Work_uiduser

	// create insert and transaction
	sqlStatement := `INSERT INTO ` + nameTableUpload + ` (work_uid,work_uiduser,work_uidfolder,work_name,work_uidwks,work_user_up,
		work_user_criacao,work_ip,work_browser,work_type,work_size,work_timespent)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`

	stmt, err := Db.Prepare(sqlStatement)

	if err != nil {

		return false, `{"status":"error","msg":"Prepare base [` + fmt.Sprintf("%s", err) + `]!"}`
	}

	if _, err := stmt.Exec(Up.Work_uid, Up.Work_uiduser, Up.Work_uidfolder, Up.Work_name, Up.Work_uidwks, Up.Work_user_up,
		Up.Work_user_criacao, Up.Work_ip, Up.Work_browser, Up.Work_type, Up.Work_size, Up.Work_timespent); err != nil {

		//log.Println(Up.Work_uid, Up.Work_uiduser, Up.Work_uidfolder, Up.Work_name, Up.Work_uidwks, Up.Work_user_up, Up.Work_user_criacao, Up.Work_ip, Up.Work_browser, Up.Work_type, Up.Work_size, Up.Work_timespent)

		return false, `{"status":"error","msg":"Exec base upload insert [` + fmt.Sprintf("%s", err) + `]!"}`
	}

	return true, `{"status":"success","msg":"ok"}`
}

func (Up *PgWorkUpload) PgUploadUpdate() (bool, string) {

	// start connect
	// Postgresql
	// local
	var Db *sql.DB

	// Db...
	if interf := pg.PgConnect(pg.DB_NAME_1, pg.DB_HOST_1, pg.DB_USER_1, pg.DB_PASSWORD_1, pg.DB_PORT_1, pg.DB_SORCE_1); interf != nil {

		Db = interf.(*sql.DB)

	} else {

		log.Println("error ao fazer connect PgUploadUpdate com Db..", interf)
		return false, "error ao fazer connect PgUploadUpdate"
	}

	nameTableUpload := "fi_" + Up.Work_uiduser

	sqlStatement := `UPDATE ` + nameTableUpload + ` SET work_user_up=$1,work_hora_up=$2,work_data_up=$3, work_ip=$4, work_browser=$5, work_size=$6, work_timespent=$7 WHERE work_name=$8 AND work_uidfolder=$9`

	stmt, err := Db.Prepare(sqlStatement)

	if err != nil {

		return false, `{"status":"error","msg":"Prepare base update [` + fmt.Sprintf("%s", err) + `]!"}`
	}

	if _, err := stmt.Exec(Up.Work_user_up, Up.Work_hora_up, Up.Work_data_up, Up.Work_ip, Up.Work_browser, Up.Work_size, Up.Work_timespent, Up.Work_name, Up.Work_uidfolder); err != nil {

		return false, `{"status":"error","msg":"Exec base update [` + fmt.Sprintf("%s", err) + `]!"}`
	}

	return true, `{"status":"success","msg":"ok"}`
}

// user authentication, the method returns user uid
// and its workspace, there are no users
// without workspaces
func PgUploadNameExist(UidUser, UidFile, work_uidfolder string) string {

	// start connect
	// Postgresql
	// local
	var Db *sql.DB

	// Db...
	if interf := pg.PgConnect(pg.DB_NAME_1, pg.DB_HOST_1, pg.DB_USER_1, pg.DB_PASSWORD_1, pg.DB_PORT_1, pg.DB_SORCE_1); interf != nil {

		Db = interf.(*sql.DB)

	} else {

		log.Println("error ao fazer connect PgUploadUpdate com Db..", interf)
		return "error ao fazer connect PgUploadUpdate"
	}

	var work_uid string

	var err error

	nameTableUpload := "fi_" + UidUser

	err = Db.QueryRow("SELECT work_uid FROM "+nameTableUpload+" WHERE work_name=$1", UidFile).Scan(&work_uid)

	switch {

	case err == sql.ErrNoRows:

		return ""

	case err != nil:

		return ""

	default:

		return work_uid
	}
}

func PgGetUploadUidFolder(UidUser string) string {

	// start connect
	// Postgresql
	// local
	var Db *sql.DB

	// Db...
	if interf := pg.PgConnect(pg.DB_NAME_1, pg.DB_HOST_1, pg.DB_USER_1, pg.DB_PASSWORD_1, pg.DB_PORT_1, pg.DB_SORCE_1); interf != nil {

		Db = interf.(*sql.DB)

	} else {

		log.Println("error ao fazer connect PgUploadUpdate com Db..", interf)
		return "error ao fazer connect PgUploadUpdate"
	}

	var workf_uid string
	var Id int

	// folder zero
	// ever 1
	Id = 1

	tableName := "fo_" + UidUser

	err := Db.QueryRow("SELECT workf_uid FROM "+tableName+" WHERE workf_id=$1", Id).Scan(&workf_uid)

	switch {

	case err == sql.ErrNoRows:

		return ""

	case err != nil:

		return ""

	default:

		//fmt.Println(workf_uid)
		return workf_uid
	}
}
