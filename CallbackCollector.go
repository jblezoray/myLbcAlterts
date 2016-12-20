package main

import (
	"time"
)

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
		adKnown, err := ac.dbAdData.IsAdKnown(*ac.curSearch, ad)
		if err != nil {
			return err
		}
		if adKnown {
			adPersisted, err := ac.dbAdData.GetAd(*ac.curSearch, ad.Id)
			if err != nil {
				return err
			}
			adPersisted.MetaData_DateSeenLast = time.Now()
			ac.dbAdData.SaveAd(*ac.curSearch, ad)

		} else {
			ad.MetaData_DateSeenFirst = time.Now()
			ac.dbAdData.SaveAd(*ac.curSearch, ad)
			ac.newAdsBySearch[*ac.curSearch] = append(ac.newAdsBySearch[*ac.curSearch], ad)
		}
	}

	return nil
}
