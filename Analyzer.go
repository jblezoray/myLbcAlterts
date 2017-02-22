package main

import (
	"fmt"
	"sort"
	"time"
)

func Analyze(adDatasOriginal []AdData) {
	adDatas := filterOnlySoldOnes(adDatasOriginal)
	sort.Sort(AdDataSortByDate(adDatas))
	printData(adDatas)
	fmt.Printf("\t\tAll Datas                   : %d ads\n", len(adDatasOriginal))
	fmt.Printf("\t\tSuspected to have been Sold : %d ads\n", len(adDatas))
	fmt.Printf("\t\tMean price                  : %d €\n", meanPrice(adDatas))
	fmt.Printf("\t\tMedian price                : %d €\n", medianPrice(adDatas))
}

func printData(adDatas []AdData) {
	printTextAbridgedHeader()
	for _, adData := range adDatas {
		printTextAbridged(adData)
	}
	printTextAbridgedFooter()
}

func filterOnlySoldOnes(adDatas []AdData) []AdData {
	adDatasCpy := []AdData{}
	h24Ago := time.Now().Local().Add(-24 * time.Hour)
	for _, adData := range adDatas {
		if !adData.MetaData_DateSeenLast.IsZero() &&
			adData.MetaData_DateSeenLast.Before(h24Ago) {
			adDatasCpy = append(adDatasCpy, adData)
		}
	}
	return adDatasCpy
}

func meanPrice(adDatas []AdData) int {
	if len(adDatas) == 0 {
		return AdDataNoPrice
	}
	total := 0
	for _, adData := range adDatas {
		if adData.Price != AdDataNoPrice {
			total += adData.Price
		}
	}
	return total / len(adDatas)
}

func medianPrice(adDatas []AdData) int {
	if len(adDatas) == 0 {
		return AdDataNoPrice
	}
	prices := []int{}
	for _, adData := range adDatas {
		if adData.Price != AdDataNoPrice {
			prices = append(prices, adData.Price)
		}
	}
	sort.Ints(prices)
	return prices[(len(prices)+1)/2]
}
