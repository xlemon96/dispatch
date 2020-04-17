package heartbeat

import (
	"log"
	"time"

	"dispatch/engine/communication"
	"dispatch/engine/worker"
	"dispatch/model/api"
)

const (
	HeartBeatInterval = 5 //心跳间隔
)

type heartbeatImpl struct {
	logger      *log.Logger
	workClient  communication.WorkerClient
	workManager worker.WorkerManager
}

func NewHeartbeatImpl(workClient communication.WorkerClient, workManager worker.WorkerManager,
	logger *log.Logger) *heartbeatImpl {
	return &heartbeatImpl{
		logger:      logger,
		workClient:  workClient,
		workManager: workManager,
	}
}

func (h *heartbeatImpl) Start() error {
	go h.doHeartbeat()
	return nil
}

func (h *heartbeatImpl) doHeartbeat() {
	for {
		workers, err := h.workManager.ListAllWorkers()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		for _, w := range workers {
			req := &api.HeartBeatRequest{
				HostIp: w.GetWorker().HostIp,
				Port:   w.GetWorker().Port,
			}
			_, err := h.workClient.Heartbeat(req)
			if err != nil {
				if err := h.workManager.UpdateWorkerState(w.GetWorker().Name, worker.WorkerStateDead); err != nil {
					//todo
				}
			} else {
				if err := h.workManager.UpdateWorkerState(w.GetWorker().Name, worker.WorkerStateRunning); err != nil {
					//todo
				}
			}
		}
		time.Sleep(time.Second * HeartBeatInterval)
	}
}
