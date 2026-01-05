package worker

import (
	"assignment_2_AP/internal/store"
	"log"
	"time"
)

type Monitor struct {
	taskStore *store.Repository[string, *store.Task]
	interval  time.Duration
	ticker    *time.Ticker
	stopChan  chan struct{}
}

// new monitoring worker
func NewMonitor(taskStore *store.Repository[string, *store.Task], interval time.Duration) *Monitor {
	return &Monitor{
		taskStore: taskStore,
		interval:  interval,
		stopChan:  make(chan struct{}),
	}
}
func (m *Monitor) Start() {
	m.ticker = time.NewTicker(m.interval)

	go m.run()
	log.Printf("Started monitoring worker (interval: %v)", m.interval)
}

// main monitoring loop
func (m *Monitor) run() {
	for {
		select {
		case <-m.ticker.C:
			m.logStats()

		case <-m.stopChan:
			m.ticker.Stop()
			log.Println("Monitoring worker stopped")
			return
		}
	}
}

// counts and logs tasks by status
func (m *Monitor) logStats() {
	tasks := m.taskStore.GetAll()
	statusCounts := map[store.TaskStatus]int{
		store.StatusPending:    0,
		store.StatusInProgress: 0,
		store.StatusDone:       0,
	}

	for _, task := range tasks {
		statusCounts[task.Status]++
	}
	log.Printf("Task Status - PENDING: %d, IN_PROGRESS: %d, DONE: %d",
		statusCounts[store.StatusPending],
		statusCounts[store.StatusInProgress],
		statusCounts[store.StatusDone])
}

func (m *Monitor) Stop() {
	log.Println("Stopping monitoring worker...")
	close(m.stopChan)
}
