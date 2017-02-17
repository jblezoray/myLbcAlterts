package main

import "testing"
import "time"
import "strings"

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

func TestParseTextDate(t *testing.T) {
	parseTextDateTesterNil(t, "Hier, 08:25", false)
	parseTextDateTesterNil(t, "Aujourd'hui, 08:08", false)
	parseTextDateTester(t, "Hier, 08:25", true, "T08:25:00+01:00")
	parseTextDateTester(t, "Aujourd'hui, 08:08", true, "T08:08:00+01:00")
	parseTextDateTester(t, "22 oct, 14:02", true, "-10-22T14:02:00+02:00")
	parseTextDateTester(t, "22 oct, 14:02", false, "-10-22T14:02:00+02:00")
	parseTextDateTester(t, "7 nov, 22:01", true, "-11-07T22:01:00+01:00")
	parseTextDateTester(t, "29 déc, 10:57", true, "-12-29T10:57:00+01:00")
	parseTextDateTester(t, "1 jan, 23:04", true, "-01-01T23:04:00+01:00")
	parseTextDateTester(t, "20 fév, 11:02", true, "-02-20T11:02:00+01:00")
}

func parseTextDateTester(t *testing.T, input string, relativeIsValid bool, expectedOutput string) {
	var ts = ParseTextDate(input, relativeIsValid)
	if ts == nil {
		t.Fatal("Expected '", expectedOutput, "' got nil")
	}
	output := ts.Format(time.RFC3339)

	if !strings.HasSuffix(output, expectedOutput) {
		t.Fatal("Expected '", expectedOutput, "' got '", output, "'")
	}
}

func parseTextDateTesterNil(t *testing.T, input string, relativeIsValid bool) {
	var output = ParseTextDate(input, relativeIsValid)
	if nil != output {
		t.Fatal("Expected nil got '", output, "'")
	}
}
