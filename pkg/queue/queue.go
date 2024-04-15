package queue

import (
	"context"
	"fmt"
	"time"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"go.uber.org/zap"
)

type MessageHandleFunc func(context.Context, WorkTicket) error

//WorkTicket ...
type WorkTicket struct {
	Type 	int
	Data	interface{}
}

//Worker ...
type Worker struct {
	ID int

	//chan chan for send and receiver workticket
	WorkItem    chan WorkTicket
	WorkerQueue chan chan WorkTicket

	//chan for quit worker
	QuitChan chan bool

	Timeout int

	// handler
	Handler MessageHandleFunc
}

//NewWorker ...
func NewWorker(id int, workerQueue chan chan WorkTicket, timeout int, handler MessageHandleFunc) Worker {
	return Worker{
		ID:				id,
		Timeout: 		timeout,
		WorkerQueue:	workerQueue,
		WorkItem:    	make(chan WorkTicket),
		QuitChan:    	make(chan bool),
		Handler:     	handler,
	}
}

//Start ...
func (w *Worker) Start() {
	go func() {
		for {
			//put to chan chan for the first time
			w.WorkerQueue <- w.WorkItem

			select {
			case workItem := <-w.WorkItem:
				go w.processWorkItem(workItem)

			case <-w.QuitChan:
				log.Bg().Debug("shut down worker")
				return
			}
		}
	}()
}

func (w *Worker) processWorkItem(ticket WorkTicket) {
	defer func() {
		if r := recover(); r != nil {
			log.Bg().Error("[process-work-item] recovering occur", zap.Any("recover", r))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(w.Timeout))
	defer cancel()

	log.WithContext(ctx).Debug("[process-work-item] process", log.Field("worker_id", w.ID), log.Field("ticket", ticket))
	w.Handler(ctx, ticket)
}

//Stop ...
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

type LocalQueue struct {
	NumberWorker int
	TimeoutPerTicket int
	QueueSize int

	//WorkerQueueVar ...
	workerQueueVar chan chan WorkTicket

	//WorkItemVar ...
	workItemVar chan WorkTicket

	//Workers ...
	workers []Worker

	// Handler ...
	Handler MessageHandleFunc
}

//SendTicket ...
func (q *LocalQueue) SendTicket(item WorkTicket) error {
	// q.workItemVar <- item
	select {
	case q.workItemVar <- item:
		return nil
	case <-time.After(time.Duration(q.TimeoutPerTicket) * time.Second):
		return fmt.Errorf("timeout")
	}
}

//StartWorkerDispatcher ...
func (q *LocalQueue) StartWorkerDispatcher() {
	//initial global var
	q.workerQueueVar = make(chan chan WorkTicket, q.QueueSize)
	q.workItemVar = make(chan WorkTicket, q.NumberWorker)

	//instantialize workers
	q.workers = make([]Worker, q.NumberWorker)

	for i := 0; i < q.NumberWorker; i++ {
		log.Bg().Debug("[start-worker-dispatcher] - start worker", log.Field("id", i+1))

		//start worker
		worker := NewWorker(i, q.workerQueueVar, q.TimeoutPerTicket, q.Handler)
		worker.Start()

		//keep worker for shutdown all workers
		q.workers[i] = worker
	}

	//start receive workItem
	go func() {
		for {
			select {
			case workItemRecv := <-q.workItemVar:
				go func() {
					worker := <-q.workerQueueVar
					worker <- workItemRecv
				}()
			}
		}
	}()
}

//StopWorkerDispatcher ...
func (c *LocalQueue) StopWorkerDispatcher() {
	for i := 0; i < len(c.workers); i++ {
		log.Bg().Debug("stop worker dispatcher", zap.Int("id", i+1))

		//stop worker
		c.workers[i].Stop()
	}
}