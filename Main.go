package main

import "os"
import "fmt"

func printUsage(progName string) {
	fmt.Println(progName + " is a scraper of a well knwon classified ads website.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("    " + progName + " configFile.json")
	fmt.Println()
}

func main() {

	// parse arguments
	args := os.Args
	if len(args) != 2 {
		printUsage(args[0])
		return
	}

	// retrieve data about previous launches / read config
	// TODO
	// query := "porsche%20924"
	config, err := ReadConfigFile(args[1])
	// config, err := ReadConfigFile("sampleConf.json")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	// build URL to scrap
	// TODO build an URL builder.
	url := "https://www.leboncoin.fr/voitures/offres/bretagne/occasions/?q=" + config.SearchTerms
	// scrap new data
	ads, err := Scraper(url)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	// it could be possible to compute a distance to a point, on the basis of this webservice :
	// $ curl "http://api-adresse.data.gouv.fr/search/?type=city&q=Carcassonne" |jq

	// print new data
	for _, ad := range ads {
		PrintTextAbridged(ad)
		//PrintText(ad)
		//PrintLineSeparator()
	}

	// build & send a mail
	// TODO
}
