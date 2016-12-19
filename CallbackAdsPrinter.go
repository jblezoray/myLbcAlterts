package main

func callbackAdsPrinter(curads []AdData, search Search) error {
	for _, ad := range curads {
		PrintTextAbridged(ad)
		//PrintText(ad)
		//PrintLineSeparator()
	}
	return nil
}
