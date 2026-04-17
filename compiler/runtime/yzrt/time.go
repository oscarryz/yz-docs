package yzrt

import gtime "time"

// Time is the built-in time singleton, accessible in Yz as `time`.
var Time = &_timeBoc{}

type _timeBoc struct{}

// Now returns the current time as a String.
func (t *_timeBoc) Now() *Thunk[String] {
	return Go(func() String {
		return NewString(gtime.Now().Format(gtime.RFC3339))
	})
}
