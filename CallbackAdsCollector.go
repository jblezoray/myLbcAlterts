package main

type CallbackAdsCollector struct {
	dbAdData    DbAdData
	adsBySearch map[Search][]AdData
}

func (ac CallbackAdsCollector) init(dbAdData DbAdData) {
	ac.dbAdData = dbAdData
	ac.adsBySearch = make(map[Search][]AdData)
}

func (ac CallbackAdsCollector) callbackAds(curads []AdData, search Search) error {

	if ac.adsBySearch[search] == nil {
		ac.adsBySearch[search] = make([]AdData, 0)
	}

	for _, ad := range curads {
		SaveAd(&ac.dbAdData, search, ad)
		ac.adsBySearch[search] = append(ac.adsBySearch[search], ad)
	}

	return nil
}
