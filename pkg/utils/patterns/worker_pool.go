package patterns

import (
	"log"
	"sync"
)

// Pool represents a pool of worker goroutines.
type Pool struct {
	numWorkers int
	tasks      chan func() error
	wg         sync.WaitGroup
}

// NewWorkerPool creates a new worker pool with a given number of workers.
func NewWorkerPool(numWorkers int) *Pool {
	return &Pool{
		numWorkers: numWorkers,
		tasks:      make(chan func() error),
	}
}

// Run starts the worker pool.
func (p *Pool) Run() {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// Submit adds a task to the pool.
func (p *Pool) Submit(t func() error) {
	p.tasks <- t
}

// Close stops accepting new tasks and waits for all workers to finish.
func (p *Pool) Close() {
	close(p.tasks)
	p.wg.Wait()
}

// worker executes tasks from the queue.
func (p *Pool) worker(id int) {
	defer p.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC recovered while processing task: %v", r)
		}
	}()
	for task := range p.tasks {
		if err := task(); err != nil {
			log.Printf("worker %d: task failed: %v", id, err)
		}
	}
}
