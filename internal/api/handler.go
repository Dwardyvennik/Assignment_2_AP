package api

import (
	"assignment_2_AP/internal/queue"
	"assignment_2_AP/internal/store"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

var taskIDCounter uint64

type Handler struct {
	taskStore *store.Repository[string, *store.Task]
	taskQueue *queue.Queue[*store.Task]
	stats     *store.Stats
}

// new API handler
func NewHandler(
	taskStore *store.Repository[string, *store.Task],
	taskQueue *queue.Queue[*store.Task],
	stats *store.Stats,
) *Handler {
	return &Handler{
		taskStore: taskStore,
		taskQueue: taskQueue,
		stats:     stats,
	}
}

// TasksHandler(w h handles /tasks endpoint
func (h *Handler) TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createTask(w, r)
	case http.MethodGet:
		h.getAllTasks(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) TaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if path == "" {
		http.Error(w, "Task ID required", http.StatusBadRequest)
		return
	}

	h.getTaskByID(w, r, path)
}

func (h *Handler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.getStats(w, r)
}

// POST /tasks
func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Payload string `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Payload == "" {
		http.Error(w, "Payload is required", http.StatusBadRequest)
		return
	}

	// task with unique id
	taskID := fmt.Sprintf("%d", atomic.AddUint64(&taskIDCounter, 1))
	task := &store.Task{
		ID:        taskID,
		Payload:   req.Payload,
		Status:    store.StatusPending,
		CreatedAt: time.Now(),
	}

	h.taskStore.Set(taskID, task)
	h.taskQueue.Enqueue(task)
	h.stats.IncrementSubmitted()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      task.ID,
		"status":  string(task.Status),
		"message": "Task created successfully",
	})
}

// Get /tasks
func (h *Handler) getAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.taskStore.GetAll()

	// Update  in  progress count
	inProgressCount := 0
	for _, task := range tasks {
		if task.Status == store.StatusInProgress {
			inProgressCount++
		}
	}
	h.stats.SetInProgress(inProgressCount)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GET /tasks/{id}
func (h *Handler) getTaskByID(w http.ResponseWriter, r *http.Request, id string) {
	task, exists := h.taskStore.Get(id)

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// GET /stats
func (h *Handler) getStats(w http.ResponseWriter, r *http.Request) {
	// update in progress count before stats
	tasks := h.taskStore.GetAll()
	inProgressCount := 0
	for _, task := range tasks {
		if task.Status == store.StatusInProgress {
			inProgressCount++
		}
	}
	h.stats.SetInProgress(inProgressCount)

	submitted, completed, inProgress := h.stats.GetStats()
	response := map[string]int{
		"submitted":   submitted,
		"completed":   completed,
		"in_progress": inProgress,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
