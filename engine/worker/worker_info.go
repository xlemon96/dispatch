package worker

import (
	"sync"

	"github.com/naivelife/dispatch/model/running"
)

const (
	WorkerStatePending = "pending"
	WorkerStateRunning = "running"
	WorkerStateLoss    = "loss"
	WorkerStateDead    = "dead"
	WorkerStateAll     = "all" // 查询状态
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
