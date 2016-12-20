package main

type CallbackCollector struct {
	dbAdData    *DbAdData
	adsBySearch map[Search][]AdData
	curSearch   *Search
}

func CallbackCollectorFactory(dbAdData DbAdData) *CallbackCollector {
	ac := &CallbackCollector{}
	ac.dbAdData = &dbAdData
	ac.adsBySearch = make(map[Search][]AdData)
	return ac
}

func (ac *CallbackCollector) callbackNewSearch(search *Search) error {
	ac.curSearch = search
	return nil
}

func (ac *CallbackCollector) callbackAds(curads []AdData) error {
	if ac.adsBySearch[*ac.curSearch] == nil {
		ac.adsBySearch[*ac.curSearch] = make([]AdData, 0)
	}

	for _, ad := range curads {
		SaveAd(ac.dbAdData, *ac.curSearch, ad)
		ac.adsBySearch[*ac.curSearch] = append(ac.adsBySearch[*ac.curSearch], ad)
	}

	return nil
}
