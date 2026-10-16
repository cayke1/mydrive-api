package jobs

import (
	"context"
	"errors"
	"log"
	"sync"
	"sync/atomic"
)

type JobQueue interface {
	Enqueue(ctx context.Context, job *Job) error
	Dequeue(ctx context.Context) (*Job, error)
	Start(ctx context.Context, workerCount int) error
	Stop()
	QueueSize() int
	ProcessedCount() int
}

type JobHandler interface {
	Handle(ctx context.Context, job *Job) error
}

type Dispatcher struct {
	jobs           chan Job
	processedCount atomic.Int64
	running        bool
	mu             sync.Mutex
	wg             sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
	handler        JobHandler
}

func NewJobDispatcher(capacity int, handler JobHandler) *Dispatcher {
	return &Dispatcher{
		jobs:    make(chan Job, capacity),
		handler: handler,
	}
}

func (d *Dispatcher) Enqueue(ctx context.Context, job *Job) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case d.jobs <- *job:
		return nil
	default:
		return errors.New("job queue is full")
	}
}

func (d *Dispatcher) Start(ctx context.Context, workerCount int) error {
	d.mu.Lock()
	if d.running {
		d.mu.Unlock()
		return errors.New("dispatcher is already running")
	}
	if workerCount <= 0 {
		d.mu.Unlock()
		return errors.New("worker count must be greater than zero")
	}

	d.ctx, d.cancel = context.WithCancel(ctx)
	d.running = true
	d.mu.Unlock()

	for i := 0; i < workerCount; i++ {
		d.wg.Add(1)
		go d.worker(d.ctx)
	}

	return nil
}

func (d *Dispatcher) worker(ctx context.Context) {
	defer d.wg.Done()
	log.Printf("[JOBS->Worker] Worker started")
	for {
		select {
		case <-ctx.Done():
			log.Printf("[JOBS->Worker] Worker shutting down")
			return
		case job, ok := <-d.jobs:
			if !ok {
				log.Printf("[JOBS->Worker] Job channel closed")
				return
			}
			d.processJob(&job)
		}
	}
}

func (d *Dispatcher) processJob(job *Job) {
	log.Printf("[JOBS->Worker] Processing job: %s", job.ID)
	for attempt := 0; attempt <= job.MaxRetries; attempt++ {
		job.RetryCount = attempt
		err := d.handler.Handle(d.ctx, job)
		if err == nil {
			d.processedCount.Add(1)
			log.Printf("[JOBS->Worker] Job processed successfully: %s, attempt %d", job.ID, attempt+1)
			return
		}
		if attempt < job.MaxRetries {
			log.Printf("[JOBS->Worker] Job failed, retrying: %s, attempt %d", job.ID, attempt+1)
			continue
		}
		log.Printf("[JOBS->Worker] Job failed after max retries: %s, attempt %d", job.ID, attempt+1)
	}
}

func (d *Dispatcher) Dequeue(ctx context.Context) (*Job, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case job, ok := <-d.jobs:
		if !ok {
			return nil, errors.New("job queue is closed")
		}
		return &job, nil
	}
}

func (d *Dispatcher) Stop() {
	d.mu.Lock()
	if !d.running {
		d.mu.Unlock()
		return
	}
	d.running = false
	d.mu.Unlock()

	log.Printf("[JOBS->Dispatcher] Stopping dispatcher")

	if d.cancel != nil {
		d.cancel()
	}

	close(d.jobs)
	d.wg.Wait()

	log.Printf("[JOBS->Dispatcher] Dispatcher stopped")
}

func (d *Dispatcher) QueueSize() int {
	return len(d.jobs)
}

func (d *Dispatcher) ProcessedCount() int {
	return int(d.processedCount.Load())
}
