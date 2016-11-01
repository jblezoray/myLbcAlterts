package main

import "fmt"

func PrintText(ad AdData) {
	fmt.Printf("----------------------------------------------------------------------------\n")
	fmt.Printf("titre    =  %s\n", ad.title)
	fmt.Printf("date     =  %s\n", ad.dateStr)
	if ad.price != NoPrice {
		fmt.Printf("price    =  %d €\n", ad.price)
	}
	fmt.Printf("categ    =  %s\n", ad.category)
	//fmt.Printf("location =  %s / %s\n", ad.locationTown, ad.locationRegion)
	fmt.Printf("thumb    =  %s\n", ad.thumbSrc)
	if ad.url != NoURL {
		fmt.Printf("url      =  %s\n", ad.url)
	}
}
