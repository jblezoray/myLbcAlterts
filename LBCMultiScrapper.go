package main

import "fmt"

type InterfaceCallbackAds interface {
	callbackAds(curads []AdData, search Search) error
}

func MultiScraper(config Configuration, dbAdData DbAdData) (map[Search][]AdData, error) {

	adsCollector := CallbackAdsCollector{}
	adsCollector.init(dbAdData)

	timer := CallbackTimer{}
	timer.init(config.TimeBetweenRequestsInSeconds)

	for _, search := range config.Searches {
		fmt.Printf("Scraping '%s'\n", search.Name)
		err := Scraper(search, adsCollector.callbackAds, callbackAdsPrinter, timer.callbackSlowMeDown)
		if err != nil {
			fmt.Print(err.Error())
			return nil, err
		}
	}

	return adsCollector.adsBySearch, nil
}
