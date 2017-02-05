package main

import "testing"
import "time"

func TestMergeWithAd(t *testing.T) {
	// Having
	today := time.Now()
	yesterday := today.Add(-time.Duration(24) * time.Hour)
	tomorrow := today.Add(time.Duration(24) * time.Hour)

	ad1 := new(AdData)
	ad1.Id = 123
	ad1.MetaData_DateSeenFirst = yesterday
	ad1.MetaData_DateSeenLast = today

	ad2 := new(AdData)
	ad2.Id = 123
	ad2.MetaData_DateSeenFirst = today
	ad2.MetaData_DateSeenLast = tomorrow

	// When
	ad2.MergeWithAd(ad1)

	// Then
	if ad2.MetaData_DateSeenFirst != yesterday || ad2.MetaData_DateSeenLast != tomorrow {
		t.Fail()
	}
}
