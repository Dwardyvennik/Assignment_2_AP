package worker

import (
	"assignment_2_AP/internal/queue"
	"assignment_2_AP/internal/store"
	"log"
	"sync"
	"time"
)

type Pool struct {
	workerCount int
	queue       *queue.Queue[*store.Task]
	taskStore   *store.Repository[string, *store.Task]
	stats       *store.Stats
	wg          sync.WaitGroup
	stopChan    chan struct{}
}

// new worker pool
func NewPool(
	workerCount int,
	queue *queue.Queue[*store.Task],
	taskStore *store.Repository[string, *store.Task],
	stats *store.Stats,
) *Pool {
	return &Pool{
		workerCount: workerCount,
		queue:       queue,
		taskStore:   taskStore,
		stats:       stats,
		stopChan:    make(chan struct{}),
	}
}

func (p *Pool) Start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	log.Printf("Started %d workers", p.workerCount)
}

// processes tasks from the queue
func (p *Pool) worker(id int) {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.queue.Dequeue():
			if !ok {
				log.Printf("Worker %d: Queue closed, exiting", id)
				return
			}

			p.processTask(id, task)
		case <-p.stopChan:
			log.Printf("Worker %d: Stop signal received, exiting", id)
			return
		}
	}
}

// habnldes task processing
func (p *Pool) processTask(workerID int, task *store.Task) {
	log.Printf("Worker %d: Processing task %s", workerID, task.ID)

	// Update in progress
	p.taskStore.Update(task.ID, func(t *store.Task) *store.Task {
		t.Status = store.StatusInProgress
		return t
	})

	processingTime := time.Duration(len(task.Payload)*100) * time.Millisecond
	if processingTime > 3*time.Second {
		processingTime = 3 * time.Second
	}
	time.Sleep(processingTime)

	// update  DONE
	p.taskStore.Update(task.ID, func(t *store.Task) *store.Task {
		t.Status = store.StatusDone
		return t
	})

	p.stats.IncrementCompleted()

	log.Printf("Worker %d: Completed task %s", workerID, task.ID)
}

func (p *Pool) Stop() {
	log.Println("Stopping worker pool...")
	close(p.stopChan)

	// wait finish with timeout
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		log.Println("All workers stopped")
	case <-time.After(5 * time.Second):
		log.Println("Worker shutdown timeout, some workers may not have finished")
	}
}
