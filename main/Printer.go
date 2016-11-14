package main

import "fmt"
import "strconv"

func PrintLineSeparator() {
	fmt.Printf("----------------------------------------------------------------------------\n")
}

func PrintRawDom(ad AdData) {
	var html, _ = ad.rawDom.Html()
	fmt.Printf("Raw dom source >>>>\n%s\n", html)
}

func PrintText(ad AdData) {
	if ad.id != NoId {
		fmt.Printf("id       =  %d\n", ad.id)
	} else {
		fmt.Printf("id       =  ?\n")
	}
	fmt.Printf("titre    =  %s\n", ad.title)
	fmt.Printf("date     =  %s\n", ad.dateStr)
	if ad.price != NoPrice {
		fmt.Printf("price    =  %d €\n", ad.price)
	} else {
		fmt.Printf("price    =  ?\n")
	}
	fmt.Printf("location =  %s / %s\n", ad.locationTown, ad.locationRegion)
	fmt.Printf("thumb    =  %s\n", ad.thumbSrc)
	if ad.url != NoURL {
		fmt.Printf("url      =  %s\n", ad.url)
	}
}

func PrintTextAbridged(ad AdData) {
	var priceStr = ""
	if ad.price != NoPrice {
		priceStr = strconv.Itoa(ad.price) + " €"
	}

	fmt.Printf("| %10.10d | %-35.35s | %15.15s | %10.10s | %20.20s | %25.25s |\n",
		ad.id, ad.title, ad.dateStr, priceStr, ad.locationRegion, ad.locationTown)
}
