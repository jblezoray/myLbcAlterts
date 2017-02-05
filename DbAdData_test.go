package main

import (
	"io/ioutil"
	"os"
	"testing"
)

// TestLoadOrCreate checks that an empty bucket is created for each search.
func TestLoadOrCreate(t *testing.T) {
	// Having
	searches := createSingletonSearches("toto")
	tmpfile := createTemporaryFile(t)
	defer deleteTemporaryFile(tmpfile)

	// when
	dbAdData, _ := LoadOrCreate(tmpfile.Name(), searches)

	// Then
	datas, _ := dbAdData.GetAllAds(searches[0])
	if datas == nil || len(datas) != 0 {
		t.Errorf("the bucket of data should be non null and empty : %v", datas)
	}

	otherSearches := createSingletonSearches("tata")
	datas, _ = dbAdData.GetAllAds(otherSearches[0])
	if datas == nil || len(datas) != 0 {
		t.Errorf("the bucket of data should be non null and empty : %v", datas)
	}
}

func TestIsAdKnown(t *testing.T) {
	// having
	searches := createSingletonSearches("toto")
	tmpfile := createTemporaryFile(t)
	defer deleteTemporaryFile(tmpfile)
	adData := createSampleAd(123)
	var dbAdData, _ = LoadOrCreate(tmpfile.Name(), searches)

	// when
	dbAdData.SaveAd(searches[0], *adData)

	// then
	if known, _ := dbAdData.IsAdKnown(searches[0], *adData); !known {
		t.Fail()
	}
	adData2 := createSampleAd(456)
	if known, _ := dbAdData.IsAdKnown(searches[0], *adData2); known {
		t.Fail()
	}
}

func TestGetAllAds(t *testing.T) {
	// having
	searches := createSingletonSearches("toto")
	tmpfile := createTemporaryFile(t)
	defer deleteTemporaryFile(tmpfile)
	var dbAdData, _ = LoadOrCreate(tmpfile.Name(), searches)
	dbAdData.SaveAd(searches[0], *createSampleAd(123))
	dbAdData.SaveAd(searches[0], *createSampleAd(456))

	// when
	var ads, _ = dbAdData.GetAllAds(searches[0])

	// then
	if len(ads) != 2 {
		t.Fail()
	}
}

func TestGetAd(t *testing.T) {
	// having
	searches := createSingletonSearches("toto")
	tmpfile := createTemporaryFile(t)
	defer deleteTemporaryFile(tmpfile)
	var dbAdData, _ = LoadOrCreate(tmpfile.Name(), searches)
	dbAdData.SaveAd(searches[0], *createSampleAd(123))
	dbAdData.SaveAd(searches[0], *createSampleAd(456))

	// when
	var ad1, _ = dbAdData.GetAd(searches[0], 123)
	var ad2, _ = dbAdData.GetAd(searches[0], 456)
	var ad3, _ = dbAdData.GetAd(searches[0], 789)

	// then
	if ad1.Id != 123 {
		t.Errorf("expected Id 123, found: %d", ad1.Id)
	}
	if ad2.Id != 456 {
		t.Errorf("expected Id 456, found: %d", ad2.Id)
	}
	if ad3 != nil {
		t.Error("expected nil value")
	}
}

func TestSaveAd(t *testing.T) {
	// having
	searches := createSingletonSearches("toto")
	tmpfile := createTemporaryFile(t)
	defer deleteTemporaryFile(tmpfile)
	var dbAdData, _ = LoadOrCreate(tmpfile.Name(), searches)

	// when
	dbAdData.SaveAd(searches[0], *createSampleAdWithTitle(123, "1"))
	dbAdData.SaveAd(searches[0], *createSampleAdWithTitle(123, "1bis"))
	dbAdData.SaveAd(searches[0], *createSampleAdWithTitle(456, "2"))

	// then
	var ad *AdData
	if ad, _ = dbAdData.GetAd(searches[0], 123); ad.Title != "1bis" {
		t.Errorf("expected title 1bis, found: %s", ad.Title)
	}
	if ad, _ = dbAdData.GetAd(searches[0], 456); ad.Title != "2" {
		t.Errorf("expected title 1bis, found: %s", ad.Title)
	}
	if ad, _ = dbAdData.GetAd(searches[0], 789); &ad == nil {
		t.Errorf("expected nil value, found: %v", ad)
	}
}

func createSingletonSearches(name string) []Search {
	search := new(Search)
	search.Name = name
	search.Terms = "tata"
	searches := make([]Search, 1)
	searches[0] = *search
	return searches
}

func createTemporaryFile(t *testing.T) *os.File {
	tmpFile, err := ioutil.TempFile("", "sample.db")
	if err != nil {
		t.Error(err)
	}
	return tmpFile
}

func deleteTemporaryFile(tmpFile *os.File) {
	os.Remove(tmpFile.Name())
}

func createSampleAd(id int) *AdData {
	return createSampleAdWithTitle(id, "")
}

func createSampleAdWithTitle(id int, title string) *AdData {
	adData := new(AdData)
	adData.Id = id
	adData.Title = title
	return adData
}
