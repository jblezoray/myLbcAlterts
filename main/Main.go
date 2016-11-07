package main

import "fmt"

func main() {
	// retrieve data about previous launches / read config
	// TODO
	query := "porsche%20924"

	// build URL to scrap
	// TODO
	url := "https://www.leboncoin.fr/voitures/offres/bretagne/occasions/?th=1&q=" + query

	// scrap new data
	ads, err := Scraper(url)
	if err != nil {
		fmt.Print(err.Error)
		return
	}

	// it could be possible to compute a distance to a point, on the basis of this webservice :
	// $ curl "http://api-adresse.data.gouv.fr/search/?type=city&q=Carcassonne" |jq

	// print new data
	for _, ad := range ads {
		PrintText(ad)
	}

	// build & send a mail
	// TODO
}
