package pipeline

import (
	"log"
	"sync"
)

type Pipeline[T any] struct {
	defsize int
	queue   []T
	mu      *sync.Mutex
}

// Builds the pipeline. Pass size as -1 or 0 to use the default size
func (p *Pipeline[T]) Build(size int) {
	s := p.defsize

	if size > 0 {
		s = size
	}

	p.queue = make([]T, 0, s)
}

// Appends data to the pipeline
func (p *Pipeline[T]) Push(d T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.queue = append(p.queue, d)
}

// Get first n items
func (p *Pipeline[T]) GetFirstOf(n int) []T {
	p.mu.Lock()
	defer p.mu.Unlock()

	queuesize := len(p.queue)
	n = min(queuesize, n)

	// if queuesize < n {
	// 	return p.queue[0:] // return the entire queue
	// }
	return p.queue[:n]
}

// Remove first n items
func (p *Pipeline[T]) EmptyUpto(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	queuesize := len(p.queue)
	n = min(queuesize, n)

	p.queue = p.queue[n:]
}

func (p *Pipeline[T]) Size() int {
	return len(p.queue)
}

// View pipeline
func (p *Pipeline[T]) View() {
	log.Println(" ****** Printing Pipeline Data ******* ")
	for _, d := range p.queue {
		log.Println(d)
	}
	log.Println(" ****** Printed Pipeline Data ******* ")
}
