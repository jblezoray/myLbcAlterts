package main

import "time"

// AdData represents an ad identified by a LBC ID.
type AdData struct {
	Id                     int
	Title                  string
	DateStr                string
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
