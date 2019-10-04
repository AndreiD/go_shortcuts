package database

import (
	"github.com/tidwall/buntdb"
	"log"
)

var db *buntdb.DB

func init() {
	database, _ := buntdb.Open(":memory:")
	db = database
}

func Set(key string, value string) {
	err := db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, nil)
		return err
	})
	if err != nil {
		log.Error(err)
	}
}

func Get(key string) string {
	var value string
	err := db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err == buntdb.ErrNotFound {
			value = ""
		} else {
			value = val
		}
		return nil
	})
	if err != nil {
		log.Error(err)
	}
	return value
}
