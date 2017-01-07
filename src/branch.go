package core

import (
	"github.com/boltdb/bolt"
	"encoding/json"
)

var bucket = []byte("branches")

type Branch struct {
	Path string
	Name string
	IsMerged bool
	IsOutdated bool
	Author string
	LastUpdated string
}

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
		return b.Put([]byte(branch.Path), encoded)
	})
	return err
}

func GetBranches(database *bolt.DB) ([]Branch, error) {
	var branches []Branch

	errView := database.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		errForEach := bucket.ForEach(func(key, value []byte) error {
			branch := Branch{}
			errJson := json.Unmarshal(value, &branch)
			if errJson != nil {
				return errJson
			}
			branches = append(branches, branch)
			return nil
		})
		return errForEach
	})

	return branches, errView
}