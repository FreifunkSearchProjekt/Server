package database

import "github.com/dgraph-io/badger"

func CheckIfUserIsAlreadyRegistered(communityID string) (bool, error) {
	db, DbOpenErr := OpenDB()
	if DbOpenErr != nil {
		return true, DbOpenErr
	}

	var used bool

	DBUseErr := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte("account/" + communityID))
		if err != nil {
			return err
		}
		if err == badger.ErrKeyNotFound {
			used = false
		} else {
			used = true
		}
		return nil
	})
	if DBUseErr != nil {
		return true, DBUseErr
	}
	return used, nil
}

func SaveNewUser(token, cID, cName string) error {
	db, DbOpenErr := OpenDB()
	if DbOpenErr != nil {
		return DbOpenErr
	}

	err := db.Update(func(txn *badger.Txn) error {
		cIDErr := txn.Set([]byte("account/"+cID+"/community_id"), []byte(cID))
		if cIDErr != nil {
			return cIDErr
		}

		cNameErr := txn.Set([]byte("account/"+cID+"/community_name"), []byte(cName))
		if cNameErr != nil {
			return cNameErr
		}

		tokenErr := txn.Set([]byte("account/"+cID+"/token"), []byte(token))
		return tokenErr
	})

	return err
}
