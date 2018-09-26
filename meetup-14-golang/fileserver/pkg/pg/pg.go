package pg

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

///
import (
	"database/sql"
	_ "github.com/lib/pq"
)

var once sync.Once

var pchan = make(chan string)

const LayoutDateLog = "2006-01-02 15:04:05"

// const LayoutDate = "2006-01-02"
// const LayoutHour = "15:04:05"

var (
	DB_HOST_1     = "localhost"
	DB_NAME_1     = "fileserver"
	DB_USER_1     = "fileserver"
	DB_PASSWORD_1 = "1234"

	DB_PORT_1  = "5432"
	DB_SSL     = "disable"
	DB_SORCE_1 = "postgres"
)

// caso tenhamos varios
// bancos em varios hosts
// diferentes
var (
	DB_HOST_2     = "192.168.0.43"
	DB_NAME_2     = "outro_banco"
	DB_USER_2     = "user"
	DB_PASSWORD_2 = "123456"

	DB_PORT_2  = "5432"
	DB_SORCE_2 = "mysql"
)

type PgStruct struct {
	Pgdb *sql.DB
}

type StatusMsg struct {
	Msg string `json:msg`
}

// cache sync.Map
type cache struct {
	mm sync.Map
	sync.Mutex
}

var (
	err    error
	PostDb PgStruct
)

var (
	pool = &cache{}
)

func init() {

	db_host_1 := os.Getenv("DB_HOST_1")

	if db_host_1 != "" {

		DB_HOST_1 = db_host_1
	}

}

// put sync.Map
func (c *cache) put(key, value interface{}) {

	c.Lock()
	defer c.Unlock()
	c.mm.Store(key, value)
}

// get sync.Map
func (c *cache) get(key interface{}) interface{} {

	c.Lock()
	defer c.Unlock()

	v, _ := c.mm.Load(key)
	return v

}

// setLoad... fn func() interface{}
func (c *cache) loadStore(key interface{}, fc func() interface{}) interface{} {

	c.Lock()
	defer c.Unlock()

	if v, ok := c.mm.Load(key); ok {
		return v
	}

	val := fc()
	c.mm.Store(key, val)
	return val

	//v, _ := c.mm.LoadOrStore(key, fc())
	//return v
}

// conectando de forma segura usando goroutine
func PgConnect(DB_NAME, DB_HOST, DB_USER, DB_PASSWORD, DB_PORT, DB_SORCE string) interface{} {

	// existe o banco
	if dbPg := pool.get(DB_NAME); dbPg != nil {

		// return objeto conexao
		return dbPg.(*sql.DB)

	} else {

		// removendo aspas..
		DB_NAME = strings.Replace(DB_NAME, `"`, "", -1)

		DBINFO := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSL)

		//log.Println(DBINFO)

		// func para ser executada
		// dentro do loadStore
		// quando duas ou mais
		// goroutine chegarem
		// neste mesmo momento
		// de fazer um Store
		fn := func() interface{} {

			PostDb.Pgdb, err = sql.Open(DB_SORCE, DBINFO)

			if err != nil {

				errordb := fmt.Sprintf("Unable to connection to database: %v\n", err)
				log.Println("error:: ", errordb)
				defer PostDb.Pgdb.Close()
				return nil
			}

			if ok2 := PostDb.Pgdb.Ping(); ok2 != nil {

				log.Println("connect error...: ", ok2)
				defer PostDb.Pgdb.Close()
				return nil
			}

			log.Println("connect return sucess:: client [" + DB_NAME + "]")
			return PostDb.Pgdb
		}

		// recebendo conexao

		// armazenando cache loadStore
		sqlDb := pool.loadStore(DB_NAME, fn)

		if sqlDb != nil {

			return sqlDb.(*sql.DB)

		} else {

			return nil
		}
	}
}

func (dbx *PgStruct) PgPing() error {

	db := dbx.Pgdb

	if err := db.Ping(); err == nil {

		return nil

	} else {

		return err
	}
}
