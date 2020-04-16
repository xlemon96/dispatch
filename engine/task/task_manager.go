package task

import (
	"dispatch/constant"
	"dispatch/storage"
	"dispatch/engine/worker"
	"dispatch/model/running"
)

type taskManager struct {
	storage       storage.Storage
	workerManager worker.WorkerManager
}

func NewTaskManager(storage storage.Storage, workerManager worker.WorkerManager) *taskManager {
	return &taskManager{
		storage:       storage,
		workerManager: workerManager,
	}
}

func (t *taskManager) SendTask(ip, port string, dag *running.DAGInstance) error {
	if err := t.storage.UpdateDAGInstanceState(dag.Id, ip, port, constant.DAGStateRunning); err != nil {
		return err
	}
	//todo，具体发送操作
	t.workerManager.WorkerLoadIncrease(ip, port)
	return nil
}
