package task

import (
	"log"
	"time"

	"github.com/naivelife/dispatch/constant"
	"github.com/naivelife/dispatch/engine/communication"
	"github.com/naivelife/dispatch/engine/worker"
	"github.com/naivelife/dispatch/model/api"
	"github.com/naivelife/dispatch/model/running"
	"github.com/naivelife/dispatch/storage"
	"github.com/naivelife/dispatch/util"
)

type taskManager struct {
	logger        *log.Logger
	storage       storage.TaskDao
	workerManager worker.WorkerManager
	workerClient  communication.WorkerClient
}

func NewTaskManager(storage storage.Storage, workerManager worker.WorkerManager,
	workerClient communication.WorkerClient, logger *log.Logger) *taskManager {
	return &taskManager{
		logger:        logger,
		storage:       storage,
		workerManager: workerManager,
		workerClient:  workerClient,
	}
}

func (t *taskManager) SendTask(ip, port string, dag *running.DAGInstance) error {
	if err := t.storage.UpdateDAGInstanceState(dag.Id, ip, port, constant.DAGStateRunning); err != nil {
		return err
	}
	request := &api.SendTaskRequest{
		HostIp:        ip,
		Port:          port,
		DagInstanceId: dag.Id,
	}
	_, err := t.workerClient.SendTask(request)
	if err != nil {
		if err := util.Retry(func() error {
			return t.storage.UpdateDAGInstanceState(dag.Id, "", "", constant.DAGStatePending)
		}, 3, time.Second); err != nil {
			return err
		}
		return err
	}
	t.workerManager.WorkerLoadIncrease(ip, port)
	return nil
}
