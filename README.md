Assignment 2 - Generic Concurrent Web Server

Assignment2/
 ├── go.mod
 ├── main.go
 ├── internal/
 │   ├── api/
 │   │   └── handler.go
 │   ├── queue/
 │   │   └── queue.go
 │   ├── worker/
 │   │   ├── pool.go
 │   │   └── monitor.go
 │   ├── store/
 │   │   ├── repository.go
 │   │   └── model.go
 └── README.md

features
Part A: Task API

POST /tasks - Submit new task
GET /tasks - List all tasks
GET /tasks/{id} - Get specific task
GET /stats - Server statistics

Part B: Concurrency & Worker Pool 

Buffered channel task queue (size: 100)
3 concurrent workers processing tasks
Thread-safe shared state with sync.RWMutex
Asynchronous task processing

Part C: Background Monitoring 

Logs task status counts every 5 seconds
Uses time.Ticker and select
Controlled via stop channel

Part D: Graceful Shutdown

Captures OS signals (SIGINT, SIGTERM)
Stops accepting new requests
Allows active tasks to finish
10-second timeout for cleanup

Part E: Generics 

Generic Repository[K, V] for type-safe storage
Generic Queue[T] for flexible task queuing
Full type parameter implementation

run
go run .
go build -o server
./server
