package stopwatch

import "time"

type Stopwatch struct {
	stopped   bool
	StartTime time.Time
	EndTime   time.Time
	Elapsed   time.Duration
}

func (self *Stopwatch) Start() {
	self.StartTime = time.Now().UTC()
}

func (self *Stopwatch) Stop() {
	self.EndTime = time.Now().UTC()
	self.stopped = true
}

func (self *Stopwatch) Duration() time.Duration {
	if !self.stopped {
		self.Stop()
	}
	self.Elapsed = self.EndTime.Sub(self.StartTime)
	return self.Elapsed
}

var defaultSW *Stopwatch

func Start() {
	defaultSW = New()
	defaultSW.Start()
}

func Stop() {
	defaultSW.Stop()
}

func Duration() time.Duration {
	return defaultSW.Duration()
}

func Seconds() float64 {
	return defaultSW.Duration().Seconds()
}

func New() *Stopwatch {
	return &Stopwatch{stopped: false}
}
