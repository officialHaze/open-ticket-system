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
		queue: make([]T, 0, s),
	}
}

type Reservoir[T any] struct {
	queue []T
}

func (r *Reservoir[T]) Fill(d []T) {
	r.queue = append(r.queue, d...)
}

func (r *Reservoir[T]) Size() int {
	return len(r.queue)
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
