package database

import (
	"errors"
	"github.com/dgraph-io/badger"
	"os"
	"path/filepath"
	"sync"
)

var DB *badger.DB
var once sync.Once

func OpenDB() (db *badger.DB, err error) {
	once.Do(func() {
		// Open the data.db file. It will be created if it doesn't exist.
		if _, StatErr := os.Stat(filepath.ToSlash("./data/")); os.IsNotExist(StatErr) {
			MkdirErr := os.MkdirAll(filepath.ToSlash("./data/"), 0666)
			if MkdirErr != nil {
				err = MkdirErr
				return
			}
		}
		opts := badger.DefaultOptions
		opts.SyncWrites = false
		opts.Dir = filepath.ToSlash("./data/")
		opts.ValueDir = filepath.ToSlash("./data/")

		expDB, DBErr := badger.Open(opts)
		if DBErr != nil {
			err = DBErr
			return
		}

		DB = expDB
	})

	if DB == nil {
		err = errors.New("missing UserDB")
		return
	}

	db = DB
	return
}
