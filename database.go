package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

const databaseFilename string = "database.db"
const openTimeout time.Duration = 1 * time.Second

type DbAdData struct {
	boltDB *bolt.DB
}

func LoadOrCreate(config Configuration) (*DbAdData, error) {

	dbAdData := DbAdData{}
	var err error

	// open or create
	options := bolt.Options{Timeout: openTimeout} // avoid indefinite wait
	dbAdData.boltDB, err = bolt.Open(config.DatabaseFilepath, 0600, &options)
	if err != nil {
		return nil, err
	}

	// create buckets if they don't exist.
	err = dbAdData.boltDB.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists(bucketId(config.SearchTerms))
		return err
	})
	return &dbAdData, err
}

func IsAdKnown(dbAdData *DbAdData, searchTerms string, ad AdData) (bool, error) {
	var known = false
	err := dbAdData.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketId(searchTerms))
		key := adKey(ad)
		v := b.Get(key)
		known = (v != nil)
		return nil
	})
	return known, err
}

func SaveAd(dbAdData *DbAdData, searchTerms string, ad AdData) error {
	err := dbAdData.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketId(searchTerms))
		adBytes, err := adMarshal(ad)
		if err != nil {
			return err
		}
		err = b.Put(adKey(ad), adBytes)
		return err
	})
	return err
}

func SaveAndClose(dbAdData *DbAdData) error {
	return dbAdData.boltDB.Close()
}

func bucketId(searchTerms string) []byte {
	return []byte(searchTerms)
}

func adKey(ad AdData) []byte {
	return []byte(strconv.Itoa(ad.Id))
}

func adMarshal(ad AdData) ([]byte, error) {
	return json.Marshal(ad)
}
func adUnmarshal(adBytes []byte) (AdData, error) {
	var ad AdData
	err := json.Unmarshal(adBytes, &ad)
	return ad, err
}
