package task

import (
	"time"

	"dispatch/constant"
	"dispatch/engine/communication"
	"dispatch/engine/worker"
	"dispatch/model/api"
	"dispatch/model/running"
	"dispatch/storage"
	"dispatch/util"
)

type taskManager struct {
	storage       storage.Storage
	workerManager worker.WorkerManager
	workerClient  communication.WorkerClient
}

func NewTaskManager(storage storage.Storage, workerManager worker.WorkerManager,
	workerClient communication.WorkerClient) *taskManager {
	return &taskManager{
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
