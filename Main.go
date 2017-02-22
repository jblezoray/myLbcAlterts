package main

import "os"
import "fmt"
import "sort"

// Those values are initialized by a LD Flag
var (
	// Version contains the current git tag
	Version string
	// Build is the build date
	Build string
)

func main() {
	// parse arguments
	args := os.Args
	analyzeMode := len(args) == 3 && args[2] == "--analyze"
	migrateDbMode := len(args) == 3 && args[2] == "--migratedb"
	updateDbMode := len(args) == 2
	if !analyzeMode && !updateDbMode && !migrateDbMode {
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
	dbAdData, err := LoadOrCreate(config.DatabaseFilepath, config.Searches)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	// run in the good mode.
	if analyzeMode {
		analyzeDb(config, dbAdData)
	} else if migrateDbMode {
		migrateDb(config, dbAdData)
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
	fmt.Println("    " + progName + " configFile.json [--analyze|--migratedb]")
	fmt.Println()
	fmt.Println("Version : " + Version)
	fmt.Println("Build   : " + Build)
	fmt.Println()
}

func analyzeDb(config Configuration, dbAdData *DbAdData) {
	fmt.Println("analyzing Database")
	for _, search := range config.Searches {
		fmt.Println("Search : ", search.Name)
		adDatas, _ := dbAdData.GetAllAds(search)
		sort.Sort(AdDataSortByDate(adDatas))
		printTextAbridgedHeader()
		for _, adData := range adDatas {
			printTextAbridged(adData)
		}
		printTextAbridgedFooter()
	}
}

func migrateDb(config Configuration, dbAdData *DbAdData) {
	fmt.Println("migrating Database")
	err := dbAdData.Migrate()
	if err != nil {
		fmt.Print(err.Error())
		return
	}
}

func updateDb(config Configuration, dbAdData *DbAdData) {

	// Scrap new data
	collector := CallbackCollectorFactory(*dbAdData)
	timer := CallbackTimerFactory(config.TimeBetweenRequestsInSeconds)
	printer := CallbackPrinterFactory()
	// it could be possible to have a callback that compute a distance to a
	// point, and enrichies the AdData. on the basis of this webservice :
	// $ curl -s "http://api-adresse.data.gouv.fr/search/?type=city&q=Carcassonne" |jq
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
