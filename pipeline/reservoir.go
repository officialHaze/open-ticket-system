package pipeline

import (
	"log"
	"ots/settings"
)

// Pass size as 0/-1 to use default size
func NewReservoir[T any](size int) *Reservoir[T] {
	s := settings.MySettings.Get_ReservoirSize()

	if size > 0 {
		s = size
	}

	return &Reservoir[T]{
		queue:      make([]T, 0, s),
		discardbin: make([]T, 0, s),
	}
}

type Reservoir[T any] struct {
	queue      []T
	discardbin []T
}

func (r *Reservoir[T]) Fill(d []T) {
	r.queue = append(r.queue, d...)
}

func (r *Reservoir[T]) Size() int {
	return len(r.queue)
}

func (r *Reservoir[T]) BinSize() int {
	return len(r.discardbin)
}

func (r *Reservoir[T]) QueueToBin() {
	if len(r.queue) <= 0 {
		return
	}

	first := r.queue[0]
	r.discardbin = append(r.discardbin, first)

	// restructure the queue
	r.queue = r.queue[1:] // ignore the first element
}

func (r *Reservoir[T]) EmptyBin() {
	r.discardbin = r.discardbin[:0]
}

// Print the contents of reservoir
func (r *Reservoir[T]) View() {
	log.Println(" ****** Printing contents of reservoir ******")
	for _, d := range r.queue {
		log.Println(d)
	}
	log.Println(" ****** Printed contents of reservoir ******")
}

// Getters
func (r *Reservoir[T]) Get_Queue() []T {
	return r.queue
}
