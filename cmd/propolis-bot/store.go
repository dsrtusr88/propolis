package main

import (
	"time"

	"github.com/dgraph-io/badger/v2"
	"gitlab.com/catastrophic/assistance/logthis"
)

func AddRelease(db *badger.DB, path, idOrError string) error {
	// else creating entry
	// the key is path + id or error message, to have unique entries and not repeat errors.
	err := db.Update(func(txn *badger.Txn) error {
		// keeping the release one week only
		e := badger.NewEntry([]byte(path+idOrError), []byte{}).WithTTL(24 * 7 * time.Hour)
		return txn.SetEntry(e)
	})
	return err
}

func IsKnown(db *badger.DB, path, idOrError string) bool {
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(path + idOrError))
		if err != nil {
			return err
		}
		_, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		logthis.Error(err, logthis.VERBOSESTEST)
		return false
	}
	return true
}
