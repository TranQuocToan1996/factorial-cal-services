package patterns

import (
	"fmt"
	"log"
	"sync"
)

// Semaphore controls concurrency using a weighted semaphore pattern.
type Semaphore struct {
	limit chan struct{}
	wg    sync.WaitGroup
}

// NewSemaphore creates a new semaphore with a concurrency limit.
func NewSemaphore(maxConcurrent int) *Semaphore {
	if maxConcurrent <= 0 {
		maxConcurrent = 1
		fmt.Println("maxConcurrent is less than 1, setting to 1")
	}
	return &Semaphore{
		limit: make(chan struct{}, maxConcurrent),
	}
}

// Submit runs the given task, blocking if the concurrency limit is reached.
func (s *Semaphore) Submit(fn func() error) {
	s.limit <- struct{}{} // Acquire
	s.wg.Add(1)

	go func() {
		defer func() {
			<-s.limit // Release
			s.wg.Done()
		}()
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC recovered while processing task: %v", r)
			}
		}()
		if err := fn(); err != nil {
			log.Printf("task failed: %v", err)
		}
	}()
}

// Wait waits for all submitted tasks to finish.
func (s *Semaphore) Wait() {
	s.wg.Wait()
}
