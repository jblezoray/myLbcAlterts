package main

import "time"

type CallbackTimer struct {
	timeBetweenRequests time.Duration
}

func (ct CallbackTimer) init(timeBetweenRequestsInSeconds int) {
	ct.timeBetweenRequests = time.Duration(timeBetweenRequestsInSeconds) * time.Second
}

func (ct CallbackTimer) callbackSlowMeDown(curads []AdData, search Search) error {
	time.Sleep(ct.timeBetweenRequests)
	return nil
}
