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

	// prepare callbacks
	collector := CallbackCollectorFactory(*dbAdData)
	timer := CallbackTimerFactory(config.TimeBetweenRequestsInSeconds)
	printer := CallbackPrinterFactory()
	// it could be possible to have a callback that compute a distance to a
	// point, and enrichies the AdData. on the basis of this webservice :
	// $ curl -s "http://api-adresse.data.gouv.fr/search/?type=city&q=Carcassonne" |jq

	// Scrap new data
	if err := Scraper(config.Searches, collector, timer, printer); err != nil {
		fmt.Print(err.Error())
		return
	}

	// build & send a mail
	fmt.Println("Sending mail")
	if err := SendAdsByMail(config, collector.newAdsBySearch); err != nil {
		fmt.Print(err.Error())
		return
	}

	fmt.Println("DONE")
}
