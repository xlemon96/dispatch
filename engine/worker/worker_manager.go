package worker

import (
	"log"
	"math/rand"
	"sync"

	"dispatch/constant"
	"dispatch/storage"
)

type workerManager struct {
	logger  *log.Logger
	workers *sync.Map
	storage storage.Storage
}

func NewWorkerManager(storage storage.Storage, logger *log.Logger) *workerManager {
	return &workerManager{
		logger:  logger,
		workers: &sync.Map{},
		storage: storage,
	}
}

func (w *workerManager) Start() error {
	if err := w.init(); err != nil {
		return err
	}
	return nil
}

func (w *workerManager) Select() *WorkerInfo {
	workers := w.listWorkersByState(WorkerStateRunning)
	if len(workers) == 0 {
		return nil
	}
	return workers[rand.Intn(len(workers))]
}

func (w *workerManager) WorkerLoadIncrease(hostIp, port string) {
	w.workers.Range(func(key, value interface{}) bool {
		if worker, ok := value.(*WorkerInfo); ok {
			if worker.worker.HostIp == hostIp && worker.worker.Port == port {
				worker.Increace()
				return false
			}
		}
		return true
	})
}

func (w *workerManager) WorkerLoadDecrease(hostIp, port string) {
	w.workers.Range(func(key, value interface{}) bool {
		if worker, ok := value.(*WorkerInfo); ok {
			if worker.worker.HostIp == hostIp && worker.worker.Port == port {
				worker.Decreace()
				return false
			}
		}
		return true
	})
}

func (w *workerManager) ListAllWorkers() ([]*WorkerInfo, error) {
	workers := make([]*WorkerInfo, 0)
	w.workers.Range(func(key, value interface{}) bool {
		if worker, ok := value.(*WorkerInfo); ok {
			workers = append(workers, worker)
		}
		return true
	})
	return workers, nil
}

func (w *workerManager) UpdateWorkerState(name, state string) error {
	if v, ok := w.workers.Load(name); ok {
		workerInfo := v.(*WorkerInfo)
		workerInfo.GetWorker().WorkerState = state
		if workerInfo.GetWorker().WorkerState == WorkerStateDead {
			//todo
		}
	}
	if err := w.storage.UpdateWorker(name, state); err != nil {
		//todo
	}
	return nil
}

func (w *workerManager) init() error {
	workers, err := w.storage.DescribeWorkers()
	if err != nil {
		return err
	}
	dags, err := w.storage.DescribeDAGInstancesByState(constant.DAGStateRunning)
	if err != nil {
		return err
	}
	for _, worker := range workers {
		workerInfo := &WorkerInfo{
			Mutex:        sync.Mutex{},
			runningTasks: 0,
			worker:       worker,
		}
		w.workers.Store(worker.Name, workerInfo)
	}
	for _, dag := range dags {
		w.WorkerLoadIncrease(dag.HostIp, dag.Port)
	}
	return nil
}

func (w *workerManager) listWorkersByState(state string) []*WorkerInfo {
	workers := make([]*WorkerInfo, 0)
	w.workers.Range(func(key, value interface{}) bool {
		if worker, ok := value.(*WorkerInfo); ok {
			if worker.GetWorker().WorkerState == state {
				workers = append(workers, worker)
			}
		}
		return true
	})
	return workers
}
