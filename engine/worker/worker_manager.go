package worker

import (
	"sync"

	"dispatch/constant"
	"dispatch/dao"
)

type workerManager struct {
	workers *sync.Map
	storage dao.Storage
}

func NewWorkerManager(storage dao.Storage) *workerManager {
	return &workerManager{
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

func (w *workerManager) Select() (*WorkerInfo, error) {
	return nil, nil
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