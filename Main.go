package main

import "os"
import "fmt"

func main() {

	// parse arguments
	args := os.Args
	analyzeMode := len(args) == 3 && args[2] == "--analyze"
	updateDbMode := len(args) == 2
	if !analyzeMode && !updateDbMode {
		printUsage(args[0])
		return
	}
	configFilePath := args[1]

	// read configuration
	config, err := ReadConfigFile(configFilePath)
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

	if analyzeMode {
		analyzeDb(config, dbAdData)
	} else if updateDbMode {
		updateDb(config, dbAdData)
	}

	fmt.Println("DONE")
}

func printUsage(progName string) {
	fmt.Println(progName + " is a scraper of a well knwon classified ads website.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Println("    " + progName + " configFile.json [--analyze]")
	fmt.Println()
}

func analyzeDb(config Configuration, dbAdData *DbAdData) {
	fmt.Println("analyzing Database")
	for _, search := range config.Searches {
		fmt.Println("Search : ", search.Name)
		adDatas, _ := dbAdData.GetAllAds(search)
		for _, adData := range adDatas {
			fmt.Println("printing", adData)
			printText(adData)
		}
	}
}

func updateDb(config Configuration, dbAdData *DbAdData) {

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
}
