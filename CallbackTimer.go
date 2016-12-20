package main

import "time"

type CallbackTimer struct {
	timeBetweenRequests time.Duration
}

func CallbackTimerFactory(timeBetweenRequestsInSeconds int) *CallbackTimer {
	ct := &CallbackTimer{}
	ct.timeBetweenRequests = time.Duration(timeBetweenRequestsInSeconds) * time.Second
	return ct
}

func (ct *CallbackTimer) callbackNewSearch(search *Search) error {
	return nil // noop
}

func (ct *CallbackTimer) callbackAds(curads []AdData) error {
	time.Sleep(ct.timeBetweenRequests)
	return nil
}
