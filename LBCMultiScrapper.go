package main

import (
	"fmt"
	"time"
)

var globalDbAdData DbAdData
var globalAdsBySearch = make(map[Search][]AdData)
var globalTimeBetweenRequests time.Duration

func MultiScraper(config Configuration, dbAdData DbAdData) (map[Search][]AdData, error) {

	globalDbAdData = dbAdData
	globalAdsBySearch = make(map[Search][]AdData)
	globalTimeBetweenRequests = time.Duration(config.TimeBetweenRequestsInSeconds) * time.Second

	for _, search := range config.Searches {
		fmt.Printf("Scraping '%s'\n", search.Name)
		// *dbAdData,
		err := Scraper(search, callbackCollectAds, callbackPrintCurentResults, callbackSlowMeDown)
		if err != nil {
			fmt.Print(err.Error())
			return nil, err
		}
	}

	return globalAdsBySearch, nil
}

func callbackCollectAds(curads []AdData, search Search) error {

	if globalAdsBySearch[search] == nil {
		globalAdsBySearch[search] = make([]AdData, 0)
	}

	for _, ad := range curads {
		SaveAd(&globalDbAdData, search, ad)
		globalAdsBySearch[search] = append(globalAdsBySearch[search], ad)
	}

	return nil
}

func callbackPrintCurentResults(curads []AdData, search Search) error {
	for _, ad := range curads {
		PrintTextAbridged(ad)
		//PrintText(ad)
		//PrintLineSeparator()
	}
	return nil
}

// call regulation
func callbackSlowMeDown(curads []AdData, search Search) error {
	time.Sleep(globalTimeBetweenRequests)
	return nil
}
