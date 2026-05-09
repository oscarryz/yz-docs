package rt

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

// Sleep waits for the specified number of seconds.
func (t *_timeBoc) Sleep(seconds Int) *Thunk[Unit] {
	return Go(func() Unit {
		gtime.Sleep(gtime.Duration(seconds.GoInt()) * gtime.Second)
		return TheUnit
	})
}
