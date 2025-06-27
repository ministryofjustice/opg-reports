package utils

import "time"

type stopwatch struct {
	startTime time.Time
	endTime   time.Time
}

func (self *stopwatch) Start() {
	self.startTime = time.Now().UTC()
}

func (self *stopwatch) Stop() (duration time.Duration) {
	self.endTime = time.Now().UTC()
	duration = self.endTime.Sub(self.startTime)
	return
}
func (self *stopwatch) Duration() (duration time.Duration) {
	duration = self.endTime.Sub(self.startTime)
	return
}

func Stopwatch() *stopwatch {
	return &stopwatch{startTime: time.Now().UTC(), endTime: time.Now().UTC()}
}
