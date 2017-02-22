package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

const openTimeout time.Duration = 1 * time.Second

// DbAdData contains a pointer to an open database.
type DbAdData struct {
	boltDB *bolt.DB
}

// LoadOrCreate opens the database file described in dbFilePath.
// For each 'Search' entity of the searches, it asserts a bucket of data
// exists and otherwise creates it.
func LoadOrCreate(dbFilePath string, searches []Search) (*DbAdData, error) {

	dbAdData := DbAdData{}
	var err error

	// open or create
	options := bolt.Options{Timeout: openTimeout} // avoid indefinite wait
	dbAdData.boltDB, err = bolt.Open(dbFilePath, 0600, &options)
	if err != nil {
		return nil, err
	}

	// create buckets if they don't exist.
	for _, search := range searches {
		err = dbAdData.boltDB.Update(func(tx *bolt.Tx) error {
			_, err = tx.CreateBucketIfNotExists(dbAdData.bucketID(search.Name))
			return err
		})
	}
	return &dbAdData, err
}

// IsAdKnown checks if an AdData is known relatively to a Search entity.
func (dbAdData *DbAdData) IsAdKnown(search Search, ad AdData) (bool, error) {
	var known = false
	err := dbAdData.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(dbAdData.bucketID(search.Name))
		if b == nil {
			return nil
		}
		key := dbAdData.adKey(ad.Id)
		v := b.Get(key)
		known = (v != nil)
		return nil
	})
	return known, err
}

// GetAllAds returns all the ads of a search.
func (dbAdData *DbAdData) GetAllAds(search Search) ([]AdData, error) {
	var ads = make([]AdData, 0)
	err := dbAdData.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(dbAdData.bucketID(search.Name))
		if b == nil {
			return nil
		}
		c := b.Cursor()
		for id, adDataRaw := c.First(); id != nil; id, adDataRaw = c.Next() {
			var adData *AdData
			var err2 error
			if adData, err2 = dbAdData.adUnmarshal(adDataRaw); err2 != nil {
				return err2
			}
			ads = append(ads, *adData)
		}
		return nil
	})
	return ads, err
}

// GetAd returns the content of an ad on the basis of its unique id.
func (dbAdData *DbAdData) GetAd(search Search, adID int) (*AdData, error) {
	var adData *AdData
	err := dbAdData.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(dbAdData.bucketID(search.Name))
		if b == nil {
			return nil
		}
		key := dbAdData.adKey(adID)
		adDataRaw := b.Get(key)
		if adDataRaw == nil {
			adData = nil
			return nil
		}
		var err2 error
		adData, err2 = dbAdData.adUnmarshal(adDataRaw)
		return err2
	})
	return adData, err
}

// SaveAd persists the data of an ad in a Search bucket.
// if the ad already exists, the data of both ads is merged.
func (dbAdData *DbAdData) SaveAd(search Search, ad AdData) error {
	err := dbAdData.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(dbAdData.bucketID(search.Name))
		if b == nil {
			return nil
		}

		previousAdData, err := dbAdData.GetAd(search, ad.Id)
		if err != nil {
			return err
		}
		if previousAdData != nil {
			migrateAd(previousAdData)
			ad.MergeWithAd(previousAdData)
		}

		var adBytes []byte
		adBytes, err = dbAdData.adMarshal(ad)
		if err != nil {
			return err
		}
		var adKey = dbAdData.adKey(ad.Id)
		err = b.Put(adKey, adBytes)
		return err
	})
	return err
}

func (dbAdData *DbAdData) Migrate() error {
	err := dbAdData.boltDB.Update(func(tx *bolt.Tx) error {
		// for each bucket
		c := tx.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			// buckets have null values.
			if v != nil {
				continue
			}
			b := tx.Bucket(k)
			if b == nil {
				return nil
			}

			cAd := b.Cursor()
			for id, adDataRaw := cAd.First(); id != nil; id, adDataRaw = cAd.Next() {
				// read
				if adDataRaw == nil {
					continue
				}
				var adData *AdData
				var err error
				if adData, err = dbAdData.adUnmarshal(adDataRaw); err != nil {
					return err
				}

				// update
				migrateAd(adData)

				// save
				var adBytes []byte
				adBytes, err = dbAdData.adMarshal(*adData)
				if err != nil {
					return err
				}
				var adKey = dbAdData.adKey(adData.Id)
				err = b.Put(adKey, adBytes)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return dbAdData.SaveAndClose()
}

func migrateAd(adData *AdData) {
	if adData.Date.IsZero() && adData.DateStr != "" {
		t := ParseTextDate(adData.DateStr, false)
		if t != nil {
			adData.Date = *t
		}
	}
}

// SaveAndClose closes all open files descriptors in an DbAdData.
func (dbAdData *DbAdData) SaveAndClose() error {
	return dbAdData.boltDB.Close()
}

func (dbAdData DbAdData) bucketID(searchTerms string) []byte {
	return []byte(searchTerms)
}

func (dbAdData DbAdData) adKey(adID int) []byte {
	return []byte(strconv.Itoa(adID))
}

func (dbAdData DbAdData) adMarshal(ad AdData) ([]byte, error) {
	return json.Marshal(ad)
}

func (dbAdData DbAdData) adUnmarshal(adBytes []byte) (*AdData, error) {
	var ad *AdData
	err := json.Unmarshal(adBytes, &ad)
	return ad, err
}
