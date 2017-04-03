package main

import (
	"regexp"
	"strconv"
	"time"
)

// AdData represents an ad identified by a LBC ID.
type AdData struct {
	Id                     int
	Title                  string
	DateStr                string
	Date                   time.Time
	UrgentFlag             bool
	Url                    string
	Price                  int
	LocationTown           string
	LocationRegion         string
	ThumbSrc               string
	MetaData_DateSeenFirst time.Time
	MetaData_DateSeenLast  time.Time
}

// AdDataNoId is the value for AdData.LbcID if the id is innexistent.
const AdDataNoId = -1

// AdDataNoPrice is the value for AdData.Price if AdData has no price.
const AdDataNoPrice = -1

// AdDataNoURL is the value for AdData.URL if AdData has no URL
const AdDataNoURL = "noUrl"

// AdDataDefaultThumb is the value for AdData.ThumbSrc if AdData has no thumb
const AdDataDefaultThumb = "http://static.leboncoin.fr/img/no-picture.png"

// MergeWithAd merges the data of two ads.
// Note that both ads must have the same LbcID, otherwise the result is undetermined.
func (adData *AdData) MergeWithAd(otherOlderAdData *AdData) {
	adData.MetaData_DateSeenFirst = otherOlderAdData.MetaData_DateSeenFirst
}

var patternA = regexp.MustCompile(`^([1-3]?[0-9]) ([a-zé]{3,4}), ([0-9]{2}):([0-9]{2})$`)
var patternB = regexp.MustCompile(`^(Hier|Aujourd'hui), ([0-9]{2}):([0-9]{2})$`)
var monthFR = [...]string{"", "jan", "fév", "mars", "avr", "mai", "juin", "jul", "août", "sept", "oct", "nov", "déc"}

// ParseTextDate parses a "textual date" as scrapped on the website.
// dateStr is the date in string format;
// relative is a flag that indicates whether relative dates should be considered as valid or not.
func ParseTextDate(dateStr string, relativeIsValid bool) *time.Time {

	var now = time.Now().Local()
	var t *time.Time
	if matches := patternA.FindStringSubmatch(dateStr); matches != nil {
		var dayInMon, _ = strconv.Atoi(matches[1])
		var mon = frenchMonthToInt(matches[2])
		if mon < 1 || mon > 12 {
			return nil
		}
		var hours, _ = strconv.Atoi(matches[3])
		var minutes, _ = strconv.Atoi(matches[4])
		var tt = time.Date(now.Year(), mon, dayInMon, hours, minutes, 0, 0, now.Location())
		if now.Before(tt) {
			// last year !
			tt = time.Date(now.Year()-1, mon, dayInMon, hours, minutes, 0, 0, now.Location())
		}
		t = &tt

	} else if matches := patternB.FindStringSubmatch(dateStr); matches != nil {
		if !relativeIsValid {
			return nil
		}
		var hours, _ = strconv.Atoi(matches[2])
		var minutes, _ = strconv.Atoi(matches[3])
		var yesterday = 0
		if matches[1] == "Hier" {
			yesterday = 1
		}
		var tt = time.Date(now.Year(), now.Month(), now.Day()-yesterday, hours, minutes, 0, 0, now.Location())
		t = &tt
	}

	return t
}

func frenchMonthToInt(threeLetterMonth string) time.Month {
	for i := range monthFR {
		if monthFR[i] == threeLetterMonth {
			return time.Month(i)
		}
	}
	return -1
}
