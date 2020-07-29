package aggregator

import (
	"math"
	"math/rand"
	"sync"
	"testing"
)

func isEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.001
}

func TestAgg(t *testing.T) {
	const length = 50
	const inputLength = 500
	a := NewAggregator(length)
	input := make([]float64, inputLength)
	for i := range input {
		input[i] = rand.Float64() * 100
		a.Add(input[i])
	}
	var expected float64
	for _, n := range input[(inputLength - length):] {
		expected += n
	}
	if a.Length() != length {
		t.Fatalf("expected len %d but got %d", length, a.Length())
	}
	if !isEqual(expected, a.Sum()) {
		t.Fatalf("expected sum %f but got %f", expected, a.Sum())
	}
	if !isEqual(expected/length, a.Avg()) {
		t.Fatalf("expected avg %f but got %f", expected/length, a.Avg())
	}
}

func TestAggThreadSafe(t *testing.T) {
	const length = 50
	const inputLength = 500
	const threads = 8
	a := NewAggregator(length)
	generate := make([]float64, 0, inputLength*threads)
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(threads)
	adding := func() {
		for i := 0; i < inputLength; i++ {
			m.Lock()
			n := rand.Float64() * 100
			a.Add(n)
			generate = append(generate, n)
			m.Unlock()
		}
		wg.Done()
	}
	for i := 0; i < threads; i++ {
		go adding()
	}
	wg.Wait()
	var expected float64
	for _, n := range generate[(inputLength*threads - length):] {
		expected += n
	}
	if a.Length() != length {
		t.Fatalf("expected len %d but got %d", length, a.Length())
	}
	if !isEqual(expected, a.Sum()) {
		t.Fatalf("expected sum %f but got %f", expected, a.Sum())
	}
	if !isEqual(expected/length, a.Avg()) {
		t.Fatalf("expected avg %f but got %f", expected/length, a.Avg())
	}
}

func TestAggThreadSafeNoLimit(t *testing.T) {
	const inputLength = 500
	const threads = 8
	a := NewAggregator(0)
	generate := make([]float64, 0, inputLength*threads)
	m := sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(threads)
	adding := func() {
		for i := 0; i < inputLength; i++ {
			m.Lock()
			n := rand.Float64() * 100
			a.Add(n)
			generate = append(generate, n)
			m.Unlock()
		}
		wg.Done()
	}
	for i := 0; i < threads; i++ {
		go adding()
	}
	wg.Wait()
	var expected float64
	for _, n := range generate {
		expected += n
	}
	if a.Length() != inputLength*threads {
		t.Fatalf("expected len %d but got %d", inputLength*threads, a.Length())
	}
	if !isEqual(expected, a.Sum()) {
		t.Fatalf("expected sum %f but got %f", expected, a.Sum())
	}
	if !isEqual(expected/(inputLength*threads), a.Avg()) {
		t.Fatalf("expected avg %f but got %f", expected/(inputLength*threads), a.Avg())
	}
}

func TestAggNoLimit(t *testing.T) {
	const inputLength = 500
	a := NewAggregator(0)
	input := make([]float64, inputLength)
	for i := range input {
		input[i] = rand.Float64() * 100
		a.Add(input[i])
	}
	var expected float64
	for _, n := range input {
		expected += n
	}
	if a.Length() != inputLength {
		t.Fatalf("expected len %d but got %d", inputLength, a.Length())
	}
	if !isEqual(expected, a.Sum()) {
		t.Fatalf("expected sum %f but got %f", expected, a.Sum())
	}
	if !isEqual(expected/inputLength, a.Avg()) {
		t.Fatalf("expected avg %f but got %f", expected/inputLength, a.Avg())
	}
}

func BenchmarkAdd(b *testing.B) {
	a := NewAggregator(500)
	for i := 0; i < b.N; i++ {
		a.Add(1)
	}
}

func BenchmarkSum(b *testing.B) {
	a := NewAggregator(500)
	for i := 0; i < b.N; i++ {
		a.Add(1)
		a.Sum()
	}
}

func BenchmarkAvg(b *testing.B) {
	a := NewAggregator(500)
	for i := 0; i < b.N; i++ {
		a.Add(1)
		a.Avg()
	}
}

func BenchmarkAddNoLimit(b *testing.B) {
	a := NewAggregator(0)
	for i := 0; i < b.N; i++ {
		a.Add(1)
	}
}

func BenchmarkSumNoLimit(b *testing.B) {
	a := NewAggregator(0)
	for i := 0; i < b.N; i++ {
		a.Add(1)
		a.Sum()
	}
}

func BenchmarkAvgNoLimit(b *testing.B) {
	a := NewAggregator(0)
	for i := 0; i < b.N; i++ {
		a.Add(1)
		a.Avg()
	}
}
