package main

// AdDataSortByDate implements the sort Interface type
type AdDataSortByDate []AdData

func (s AdDataSortByDate) Len() int {
	return len(s)
}

func (s AdDataSortByDate) Less(i, j int) bool {
	return s[i].Date.Before(s[j].Date)
}

func (s AdDataSortByDate) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
