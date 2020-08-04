package aggregator

import (
	"container/list"
	"sync"
)

type Aggregator interface {
	Add(v float64) (float64, bool)
	Sum() float64
	Length() int
	Avg() float64
}

type aggNoLimit struct {
	sum   float64
	count int
	m     sync.Mutex
}

type agg struct {
	l   *list.List
	max int
	aggNoLimit
}

func NewAggregator(max int) Aggregator {
	if max > 0 {
		return &agg{
			l:   list.New(),
			max: max,
		}
	}
	return &aggNoLimit{}
}

func (a *agg) Add(v float64) (val float64, evicted bool) {
	a.m.Lock()
	a.sum += v
	a.l.PushFront(v)
	if a.l.Len() > a.max {
		p := a.l.Back()
		val = p.Value.(float64)
		evicted = true
		a.sum -= val
		a.l.Remove(p)
	}
	a.m.Unlock()
	return
}

func (a *agg) Length() int {
	a.m.Lock()
	ret := a.l.Len()
	a.m.Unlock()
	return ret
}

func (a *agg) Avg() float64 {
	a.m.Lock()
	ret := 0.0
	if a.l.Len() > 0 {
		ret = a.sum / float64(a.l.Len())
	}
	a.m.Unlock()
	return ret
}

func (a *aggNoLimit) Add(v float64) (float64, bool) {
	a.m.Lock()
	a.sum += v
	a.count++
	a.m.Unlock()
	return 0, false
}

func (a *aggNoLimit) Sum() float64 {
	a.m.Lock()
	s := a.sum
	a.m.Unlock()
	return s
}

func (a *aggNoLimit) Length() int {
	a.m.Lock()
	c := a.count
	a.m.Unlock()
	return c
}

func (a *aggNoLimit) Avg() float64 {
	a.m.Lock()
	ret := 0.0
	if a.count > 0 {
		ret = a.sum / float64(a.count)
	}
	a.m.Unlock()
	return ret
}
