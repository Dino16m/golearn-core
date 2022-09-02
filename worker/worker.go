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

func NewQueue(namespace string, concurrency uint, redisPool *redis.Pool) *Queue {
	return &Queue{
		namespace: namespace,
		pool:      work.NewWorkerPool(Context{}, concurrency, namespace, redisPool),
		enqueuer:  work.NewEnqueuer(namespace, redisPool),
	}
}

func (q *Queue) Start(ctx context.Context) {
	go func() {
		q.pool.Start()
		<-ctx.Done()
		q.pool.Stop()
	}()
}

func (q *Queue) Stop() {
	q.pool.Stop()
}

func (q *Queue) Enqueuer() *work.Enqueuer {
	return q.enqueuer
}

func (q *Queue) DispatchJob(jobName string, args map[string]any) {
	q.enqueuer.Enqueue(jobName, args)
}

func (q *Queue) RegisterHandler(jobName string, handler Handler) {
	q.pool.Job(jobName, func(baseJob *work.Job) error {
		job := Job{
			GocraftJob: baseJob,
			Args:       baseJob.Args,
		}
		return handler(job)
	})
}
