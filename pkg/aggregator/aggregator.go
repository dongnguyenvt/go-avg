package aggregator

import (
	"container/list"
	"sync"
)

type Aggregator interface {
	Add(v float64)
	Sum() float64
	Length() int
	Avg() float64
}

type agg struct {
	l     *list.List
	max   int
	sum   float64
	count int
	m     sync.Mutex
}

func NewAggregator(max int) Aggregator {
	return &agg{
		l:   list.New(),
		max: max,
	}
}

func (a *agg) Add(v float64) {
	a.m.Lock()
	defer a.m.Unlock()
	a.sum += v
	if a.max == 0 {
		a.count++
		return
	}
	a.l.PushFront(v)
	if a.l.Len() > a.max {
		p := a.l.Back()
		a.sum -= p.Value.(float64)
		a.l.Remove(p)
	}
}

func (a *agg) Sum() float64 {
	a.m.Lock()
	defer a.m.Unlock()
	return a.sum
}

func (a *agg) length() int {
	if a.max == 0 {
		return a.count
	}
	return a.l.Len()
}

func (a *agg) Length() int {
	a.m.Lock()
	defer a.m.Unlock()
	return a.length()
}

func (a *agg) Avg() float64 {
	a.m.Lock()
	defer a.m.Unlock()
	return a.sum / float64(a.length())
}
