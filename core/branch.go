package core

import (
	"encoding/json"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var bucket = []byte("branches")

type Branch struct {
	Repository  string
	Name        string
	FullPath    string
	IsMerged    bool
	IsOutdated  bool
	Author      string
	LastUpdated time.Time
}

func InitBranchesBucket(database *bolt.DB) {
	database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
}

// Saves branch to a storage
func (branch *Branch) Save(database *bolt.DB) error {
	err := database.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		encoded, err := json.Marshal(branch)
		if err != nil {
			return err
		}
		return b.Put([]byte(branch.FullPath), encoded)
	})
	return err
}

// Gets list of branches from a storage
func GetBranches(database *bolt.DB) ([]Branch, error) {
	var branches = []Branch{}

	InitBranchesBucket(database)

	errView := database.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(bucket)
		errForEach := bucket.ForEach(func(key, value []byte) error {
			branch := Branch{}

			if errJson := json.Unmarshal(value, &branch); errJson != nil {
				return errJson
			}
			branches = append(branches, branch)
			return nil
		})
		return errForEach
	})

	return branches, errView
}

// Cleans branches in a storage
func CleanBranches(database *bolt.DB) error {
	err := database.Update(func(tx *bolt.Tx) error {
		if errDelete := tx.DeleteBucket(bucket); errDelete != nil {
			return errDelete
		}
		return nil
	})
	return err
}
