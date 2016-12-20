package main

type CallbackCollector struct {
	dbAdData       *DbAdData
	newAdsBySearch map[Search][]AdData
	curSearch      *Search
}

func CallbackCollectorFactory(dbAdData DbAdData) *CallbackCollector {
	ac := &CallbackCollector{}
	ac.dbAdData = &dbAdData
	ac.newAdsBySearch = make(map[Search][]AdData)
	return ac
}

func (ac *CallbackCollector) callbackNewSearch(search *Search) error {
	ac.curSearch = search
	return nil
}

func (ac *CallbackCollector) callbackAds(curads []AdData) error {
	if ac.newAdsBySearch[*ac.curSearch] == nil {
		ac.newAdsBySearch[*ac.curSearch] = make([]AdData, 0)
	}

	for _, ad := range curads {
		ac.dbAdData.SaveAd(*ac.curSearch, ad)
		ac.newAdsBySearch[*ac.curSearch] = append(ac.newAdsBySearch[*ac.curSearch], ad)
	}

	return nil
}
