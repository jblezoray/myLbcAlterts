package main

import "fmt"
import "strconv"

func PrintLineSeparator() {
	fmt.Printf("----------------------------------------------------------------------------\n")
}

func PrintText(ad AdData) {
	if ad.Id != NoId {
		fmt.Printf("id       =  %d\n", ad.Id)
	} else {
		fmt.Printf("id       =  ?\n")
	}
	fmt.Printf("titre    =  %s\n", ad.Title)
	fmt.Printf("date     =  %s\n", ad.DateStr)
	if ad.Price != NoPrice {
		fmt.Printf("price    =  %d €\n", ad.Price)
	} else {
		fmt.Printf("price    =  ?\n")
	}
	fmt.Printf("location =  %s / %s\n", ad.LocationTown, ad.LocationRegion)
	fmt.Printf("thumb    =  %s\n", ad.ThumbSrc)
	if ad.Url != NoURL {
		fmt.Printf("url      =  %s\n", ad.Url)
	}
}

func PrintTextAbridged(ad AdData) {
	var priceStr = ""
	if ad.Price != NoPrice {
		priceStr = strconv.Itoa(ad.Price) + " €"
	}

	var urgent = ""
	if ad.UrgentFlag {
		urgent = "X"
	}

	fmt.Printf("| %10.10d | %-35.35s | %15.15s | %10.10s | %1s | %20.20s | %25.25s |\n",
		ad.Id, ad.Title, ad.DateStr, priceStr, urgent, ad.LocationRegion, ad.LocationTown)
}
