package main

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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

const NoId = -1
const NoPrice = -1
const NoURL = "noUrl"
const DefaultThumb = "http://static.leboncoin.fr/img/no-picture.png"

func parseAd(s *goquery.Selection) AdData {

	// parse title
	var title = s.Find("h2.item_title").Text()
	title = strings.TrimSpace(title)

	// parse date
	var dateStr = s.Find("aside.item_absolute p.item_supp").Text()
	var urgentFlag = strings.Count(dateStr, "Urgent") == 1
	dateStr = strings.Replace(dateStr, "Urgent", "", -1)
	dateStr = strings.TrimSpace(dateStr)

	// parse url
	url, exists := s.Find("a.list_item").Attr("href")
	if !exists {
		url = NoURL
	} else {
		url = "http:" + url
	}

	// parse id
	var id = NoId
	idStr, exists := s.Find("div.saveAd").Attr("data-savead-id")
	if exists {
		var err error
		id, err = strconv.Atoi(idStr)
		if err != nil {
			id = NoId
		}
	}

	// parse thumbSrc
	thumbSrc, exists := s.Find("span.lazyload").Attr("data-imgsrc")
	if !exists {
		thumbSrc = DefaultThumb
	} else {
		thumbSrc = "http:" + thumbSrc
	}

	// parse price
	var priceStr = s.Find("h3.item_price").Text()
	priceStr = strings.Replace(priceStr, "€", "", -1)
	priceStr = strings.TrimSpace(priceStr)
	priceStr = strings.Replace(priceStr, " ", "", -1)
	priceInt, err := strconv.Atoi(priceStr)
	if err != nil {
		priceInt = NoPrice
	}

	// parse category
	var category = s.Find("section.item_infos p:first").Text()
	category = strings.TrimSpace(category)

	// parse locationTown AND locationRegion :
	var rawLocation = s.Find("section.item_infos p.item_supp:nth-of-type(2n)").Text()
	var splitedLocation = strings.Split(rawLocation, "/")
	var locationTown, locationRegion string
	if len(splitedLocation) == 2 {
		locationTown = strings.TrimSpace(splitedLocation[0])
		locationRegion = strings.TrimSpace(splitedLocation[1])
	}

	adData := AdData{
		Id:             id,
		Title:          title,
		DateStr:        dateStr,
		UrgentFlag:     urgentFlag,
		Url:            url,
		Price:          priceInt,
		LocationTown:   locationTown,
		LocationRegion: locationRegion,
		ThumbSrc:       thumbSrc,
	}
	return adData
}

func scraperSinglePage(url string) ([]AdData, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, errors.New("cannot query the URL" + url)
	}

	var ads = make([]AdData, 0)
	selection := doc.Find("section.tabsContent li")
	selection.Each(func(i int, sel *goquery.Selection) {
		//debugPrintRawDom(sel)
		currentAd := parseAd(sel)
		ads = append(ads, currentAd)
	})

	return ads, nil
}

// func debugPrintRawDom(rawDom *goquery.Selection) {
// 	var html, _ = rawDom.Html()
// 	fmt.Printf("Raw dom source >>>>\n%s\n", html)
// }

type InterfaceCallbackAds interface {
	callbackNewSearch(search *Search) error
	callbackAds(curads []AdData) error
}

func Scraper(
	searches []Search,
	callbacks ...InterfaceCallbackAds) error {

	for _, search := range searches {

		for _, callback := range callbacks {
			err := callback.callbackNewSearch(&search)
			if err != nil {
				return err
			}
		}

		var url = "https://www.leboncoin.fr/" + search.Terms

		for i := 1; i < 10; i++ {
			// the "o" parameter in the page indicates the number of the page
			curads, err := scraperSinglePage(url + "&o=" + strconv.Itoa(i))
			if err != nil {
				return err
			}

			for _, callback := range callbacks {
				err = callback.callbackAds(curads)
				if err != nil {
					return err
				}
			}

			// no ad on the page
			if len(curads) == 0 {
				break
			}
		}
	}

	return nil
}
