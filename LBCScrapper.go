package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AdData struct {
	id             int
	title          string
	dateStr        string
	url            string
	price          int
	locationTown   string
	locationRegion string
	thumbSrc       string
	rawDom         *goquery.Selection
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

	return AdData{
		id:             id,
		title:          title,
		dateStr:        dateStr,
		url:            url,
		price:          priceInt,
		locationTown:   locationTown,
		locationRegion: locationRegion,
		thumbSrc:       thumbSrc,
		rawDom:         s,
	}
}

func scraperSinglePage(url string) ([]AdData, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, errors.New("cannot query the URL" + url)
	}

	fmt.Println("Scrapping ", url)

	var ads = make([]AdData, 0)
	selection := doc.Find("section.tabsContent li")
	selection.Each(func(i int, sel *goquery.Selection) {
		currentAd := parseAd(sel)
		ads = append(ads, currentAd)
	})

	return ads, nil
}

func Scraper(url string) ([]AdData, error) {

	var allAds = make([]AdData, 0)

	for i := 1; i < 10; i++ {
		// the "o" parameter in the page indicates the number of the page
		curads, err := scraperSinglePage(url + "&o=" + strconv.Itoa(i))
		if err != nil {
			return nil, err
		}

		// no more ads on the page : stop here.
		if len(curads) == 0 {
			break
		}

		// strange syntax ... see http://stackoverflow.com/a/16248257
		allAds = append(allAds, curads...)
	}

	return allAds, nil
}
