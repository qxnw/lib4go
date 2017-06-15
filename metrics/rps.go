package metrics

import "sync"

// RPS count events to produce exponentially-weighted moving average rates
// at one-, five-, and fifteen-minutes and a mean rate.
type RPS interface {
	Mark(int32)
	M1() int32
	M5() int32
	M15() int32
}

// GetOrRegisterRps returns an existing Meter or constructs and registers a
// new StandardMeter.
func GetOrRegisterRps(name string, r Registry) RPS {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, newStandardRPS).(RPS)
}

// StandardRPS is the standard implementation of a Meter.
type StandardRPS struct {
	lock sync.RWMutex
	m1   *RPSC
	m5   *RPSC
	m15  *RPSC
}

func newStandardRPS() *StandardRPS {
	return &StandardRPS{
		m1:  NewRPSC(60, 3600),
		m5:  NewRPSC(300, 3600),
		m15: NewRPSC(900, 3600),
	}
}

func (s *StandardRPS) Mark(i int32) {
	s.lock.Lock()
	s.m1.Mark(i)
	s.m5.Mark(i)
	s.m15.Mark(i)
	s.lock.Unlock()
}
func (s *StandardRPS) M1() int32 {
	return s.m1.counter
}
func (s *StandardRPS) M5() int32 {
	return s.m5.counter
}
func (s *StandardRPS) M15() int32 {
	return s.m15.counter
}
