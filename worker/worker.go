package worker

import (
	"context"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type Job struct {
	Args       map[string]any
	GocraftJob *work.Job
}

type Handler func(Job) error

type Context struct {
}

type Queue struct {
	pool      *work.WorkerPool
	namespace string
	enqueuer  *work.Enqueuer
}

// NewQueue creates a new Queue struct
func NewQueue(namespace string, concurrency uint, redisPool *redis.Pool) *Queue {
	return &Queue{
		namespace: namespace,
		pool:      work.NewWorkerPool(Context{}, concurrency, namespace, redisPool),
		enqueuer:  work.NewEnqueuer(namespace, redisPool),
	}
}

// Start starts the queue worker,
// the queue worker automatically stops if there is an interrupt
func (q *Queue) Start(ctx context.Context) {
	go func() {
		q.pool.Start()
		<-ctx.Done()
		q.pool.Stop()
	}()
}

// Stop stops the queue worker
func (q *Queue) Stop() {
	q.pool.Stop()
}

// Get the underlying gocraft enqueuer
// in cases where you want to perform more complex actions
func (q *Queue) Enqueuer() *work.Enqueuer {
	return q.enqueuer
}

// DispatchJob dispatches a jobname with an argument to the queue worker
func (q *Queue) DispatchJob(jobName string, args map[string]any) {
	q.enqueuer.Enqueue(jobName, args)
}

// RegisterHandler registers a handler function for a specific jobName
// the registration of handlers and dispatching of jobs can be done
// in any order
func (q *Queue) RegisterHandler(jobName string, handler Handler) {
	q.pool.Job(jobName, func(baseJob *work.Job) error {
		job := Job{
			GocraftJob: baseJob,
			Args:       baseJob.Args,
		}
		return handler(job)
	})
}
