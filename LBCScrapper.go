package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AdData struct {
	title          string
	dateStr        string
	url            string
	price          int
	locationTown   string
	locationRegion string
	thumbSrc       string
	rawDom         *goquery.Selection
}

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

func Scraper(url string) ([]AdData, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, errors.New("cannot query the URL" + url)
	}

	var ads = make([]AdData, 0)

	selection := doc.Find("section.tabsContent li")
	selection.Each(func(i int, sel *goquery.Selection) {
		//fmt.Printf("%d\n", i)
		currentAd := parseAd(sel)
		ads = append(ads, currentAd)
	})

	return ads, nil
}
