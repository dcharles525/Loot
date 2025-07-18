package data

import (
	"log"

	"loot/structs"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	bolt "go.etcd.io/bbolt"
)

type errMsg struct{ error }

const dbFile = "commands.db"
const bucketName = "commands"

func openDB() *bolt.DB {
	db, dbError := bolt.Open(dbFile, 0600, nil)
	if dbError != nil {
		log.Fatal(dbError)
	}
	return db
}

func closeDB(db *bolt.DB) {
	db.Close()
}

func GetAllItems() ([]list.Item, error) {
	db := openDB()
	defer closeDB(db)
	items := []list.Item{}
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket != nil {
			bucket.ForEach(func(key, value []byte) error {
				items = append(
					items,
					structs.Item{
						CommandTitle: string(key),
						CommandDesc:  string(value),
					},
				)
				return nil
			})
		}
		return nil
	})
	return items, nil
}

func CreateCommand(title string, command string) tea.Cmd {
	db := openDB()
	updateError := db.Update(func(tx *bolt.Tx) error {
		bucket, bucketError := tx.CreateBucketIfNotExists([]byte(bucketName))
		if bucketError != nil {
			return bucketError
		}
		saveError := bucket.Put([]byte(title), []byte(command))
		if saveError != nil {
			return saveError
		}
		return nil
	})
	closeDB(db)

	return func() tea.Msg {
		if updateError != nil {
			return errMsg{updateError}
		}
		return structs.UpdateCommandMsg{}
	}
}

func DeleteCommand(command string) tea.Cmd {
	db := openDB()
	updateError := db.Update(func(tx *bolt.Tx) error {
		bucket, bucketError := tx.CreateBucketIfNotExists([]byte(bucketName))
		if bucketError != nil {
			return bucketError
		}
		deleteError := bucket.Delete([]byte(command))
		if deleteError != nil {
			return deleteError
		}
		return nil
	})
	closeDB(db)

	return func() tea.Msg {
		if updateError != nil {
			return errMsg{updateError}
		}
		return structs.UpdateCommandMsg{}
	}
}
