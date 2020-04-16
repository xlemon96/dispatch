package worker

import (
	"sync"

	"dispatch/model/running"
)

type WorkerInfo struct {
	sync.Mutex
	runningTasks int32
	worker       *running.Worker
}

func (w *WorkerInfo) GetWorker() *running.Worker {
	return w.worker
}

func (w *WorkerInfo) Increace() {
	w.Lock()
	defer w.Unlock()
	w.runningTasks++
}

func (w *WorkerInfo) Decreace() {
	w.Lock()
	defer w.Unlock()
	w.runningTasks--
}
