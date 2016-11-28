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

	// read configuration
	config, err := ReadConfigFile(args[1])
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	// retrieve data about previous launches
	dbAdData, err := LoadOrCreate(config)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	// build URL to scrap & scrap new data
	// TODO build an URL builder.
	url := "https://www.leboncoin.fr/" + config.SearchTerms
	ads, err := Scraper(url, *dbAdData, config.SearchTerms)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	// it could be possible to compute a distance to a point, on the basis of this webservice :
	// $ curl "http://api-adresse.data.gouv.fr/search/?type=city&q=Carcassonne" |jq

	// print data
	for _, ad := range ads {
		PrintTextAbridged(ad)
		//PrintText(ad)
		//PrintLineSeparator()
	}

	if len(ads) == 0 {
		fmt.Println("No new Data")
	} else {
		// build & send a mail
		fmt.Println("Sending mail")
		err = SendAdsByMail(config.SMTPUser, config.MailFrom, config.MailTo, ads)
		if err != nil {
			fmt.Print(err.Error())
			return
		}
	}

	fmt.Println("DONE")
}
