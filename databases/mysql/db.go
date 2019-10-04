package mysql

import (
	_ "github.com/go-sql-driver/mysql" // needed by sqlx
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

// Db is the database reference
var Db *sqlx.DB

// InitMySQLDatabase inits the database
func InitMySQLDatabase(connectionURI string) error {
	var err error
	if Db != nil {
		err = Db.Ping()
		if err == nil {
			return nil
		}
		return err
	}

	Db, err = sqlx.Open("mysql", connectionURI)
	if err != nil {
		return err
	}

	start := time.Now()
	err = Db.Ping()
	if err != nil {
		return err
	}
	log.Printf("connection to the database OK. (%v)", time.Since(start))
	return nil

}

// CleanDB prepare the db for a new instance
func CleanDB() {
	// clean balances table
	Db.MustExec(`TRUNCATE TABLE balances`)

}

//-----------