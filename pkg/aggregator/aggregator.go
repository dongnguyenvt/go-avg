package aggregator

import (
	"container/list"
)

type agg struct {
	l     *list.List
	max   int
	sum   float64
	count int
}

func NewAggregator(max int) agg {
	return agg{
		l:   list.New(),
		max: max,
	}
}

func (a *agg) Add(v float64) {
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
	return a.sum
}

func (a *agg) Len() int {
	if a.max == 0 {
		return a.count
	}
	return a.l.Len()
}

func (a *agg) Avg() float64 {
	return a.sum / float64(a.Len())
}
