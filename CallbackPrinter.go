package main

import (
	"fmt"
	"strconv"
)

type CallbackPrinter struct {
	// noop
}

func CallbackPrinterFactory() *CallbackPrinter {
	return &CallbackPrinter{}
}

func (cp *CallbackPrinter) callbackNewSearch(search *Search) error {
	printLineSeparator()
	fmt.Printf("Scraping '%s'\n", search.Name)
	return nil
}

func (cp *CallbackPrinter) callbackAds(curads []AdData) error {
	for _, ad := range curads {
		printTextAbridged(ad)
		//printText(ad)
		//printLineSeparator()
	}
	return nil
}

func printLineSeparator() {
	fmt.Printf("----------------------------------------------------------------------------\n")
}

func printText(ad AdData) {
	if ad.Id != AdDataNoId {
		fmt.Printf("id       =  %d\n", ad.Id)
	} else {
		fmt.Printf("id       =  ?\n")
	}
	fmt.Printf("titre    =  %s\n", ad.Title)
	fmt.Printf("date     =  %s\n", ad.DateStr)
	if ad.Price != AdDataNoPrice {
		fmt.Printf("price    =  %d €\n", ad.Price)
	} else {
		fmt.Printf("price    =  ?\n")
	}
	fmt.Printf("location =  %s / %s\n", ad.LocationTown, ad.LocationRegion)
	fmt.Printf("thumb    =  %s\n", ad.ThumbSrc)
	if ad.Url != AdDataNoURL {
		fmt.Printf("url      =  %s\n", ad.Url)
	}
}

func printTextAbridged(ad AdData) {
	var priceStr = ""
	if ad.Price != AdDataNoPrice {
		priceStr = strconv.Itoa(ad.Price) + " €"
	}

	var urgent = ""
	if ad.UrgentFlag {
		urgent = "X"
	}

	fmt.Printf("| %10.10d | %-35.35s | %15.15s | %10.10s | %1s | %20.20s | %25.25s |\n",
		ad.Id, ad.Title, ad.DateStr, priceStr, urgent, ad.LocationRegion, ad.LocationTown)
}
