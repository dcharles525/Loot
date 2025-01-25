package main

import (
	"log"
	
	structs "loot/structs"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	bolt "go.etcd.io/bbolt"
)

type errMsg struct{ error }

func openDB()*bolt.DB{
	db, dbError := bolt.Open("commands.db", 0600, nil)

	if dbError != nil {
		log.Fatal(dbError)
	}

	return db
}

func closeDB(db *bolt.DB){
	db.Close()
}

func GetAllItems()[]list.Item{
	db := openDB()
	items := []list.Item{}

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("commands"))
	 
		if bucket != nil {
			bucket.ForEach(func(key, value []byte) error {
				items = append(
					items, 
					structs.Item{
						CommandTitle: string(key), 
						CommandDesc: string(value),
					},
				)
				return nil
			})
		}

		return nil
	})
	
	closeDB(db)

	return items
}

func CreateCommand(title string, command string) tea.Cmd {
	db := openDB()

	updateError := db.Update(func(tx *bolt.Tx) error {
		bucket, bucketError := tx.CreateBucketIfNotExists([]byte("commands"))
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
		bucket, bucketError := tx.CreateBucketIfNotExists([]byte("commands"))
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
