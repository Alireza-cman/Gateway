package database

import (
	"errors"
	"fmt"

	bolt "github.com/boltdb/bolt"
)

var DB *bolt.DB

// func init() {
// 	var err error
// 	DB, err = bolt.Open("bolt.db", 0644, nil)
// 	if err != nil {
// 		log.Println(err)
// 		panic(err)
// 	}

// }
func SetupDB() (*bolt.DB, error) {
	db, err := bolt.Open("test.db", 0644, nil)
	if err != nil {
		return nil, fmt.Errorf("could not open db, %v", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("DEVICE"))
		if err != nil {
			return fmt.Errorf("could not create device bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not set up buckets, %v", err)
	}
	fmt.Println("DB Setup Done")
	return db, nil
}

func StoreData(BucketName string, key []byte, value []byte) error {
	// store some data
	err := DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
func FetchData(BucketName string, key []byte, value []byte) ([]byte, error) {

	err := DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		if bucket == nil {
			return errors.New("Bucket is not found!")
		}

		val := bucket.Get(key)
		fmt.Println(string(val))

		return nil
	})
	return nil, err

}
