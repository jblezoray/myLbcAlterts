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
	for _, search := range config.Searches {
		err = dbAdData.boltDB.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucketIfNotExists(dbAdData.bucketID(search.Name))
			return err
		})
	}
	return &dbAdData, err
}

func (dbAdData *DbAdData) IsAdKnown(search Search, ad AdData) (bool, error) {
	var known = false
	err := dbAdData.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(dbAdData.bucketID(search.Name))
		key := dbAdData.adKey(ad)
		v := b.Get(key)
		known = (v != nil)
		return nil
	})
	return known, err
}

func (dbAdData *DbAdData) SaveAd(search Search, ad AdData) error {
	err := dbAdData.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(dbAdData.bucketID(search.Name))
		adBytes, err := dbAdData.adMarshal(ad)
		if err != nil {
			return err
		}
		err = b.Put(dbAdData.adKey(ad), adBytes)
		return err
	})
	return err
}

func (dbAdData *DbAdData) SaveAndClose() error {
	return dbAdData.boltDB.Close()
}

func (dbAdData DbAdData) bucketID(searchTerms string) []byte {
	return []byte(searchTerms)
}

func (dbAdData DbAdData) adKey(ad AdData) []byte {
	return []byte(strconv.Itoa(ad.Id))
}

func (dbAdData DbAdData) adMarshal(ad AdData) ([]byte, error) {
	return json.Marshal(ad)
}
func (dbAdData DbAdData) adUnmarshal(adBytes []byte) (AdData, error) {
	var ad AdData
	err := json.Unmarshal(adBytes, &ad)
	return ad, err
}
