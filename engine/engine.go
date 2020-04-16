package engine

import (
	"fmt"
	"time"

	"dispatch/engine/dispatch"
	"dispatch/engine/task"
	"dispatch/engine/worker"
	"dispatch/model/running"
	"dispatch/storage"
)

type Engine struct {
	dispatch      dispatch.Dispatch
	workerManager worker.WorkerManager
	taskManager   task.TaskManager
}

func NewEngine(dao storage.Storage) *Engine {
	Engine := &Engine{
		dispatch:      dispatch.NewDispatch(dao),
		workerManager: worker.NewWorkerManager(dao),
	}
	Engine.taskManager = task.NewTaskManager(dao, Engine.workerManager)
	return Engine
}

func (e *Engine) Start() error {
	if err := e.workerManager.Start(); err != nil {
		return err
	}
	if err := e.dispatch.Start(); err != nil {
		return err
	}
	go e.startDispatch()
	return nil
}

func (e *Engine) startDispatch() {
	todoDags := e.dispatch.GetTodoDag()
	for {
		select {
		case dagBag := <-todoDags:
			if len(dagBag.GetDagInstances()) == 0 {
				continue
			}
			for _, dag := range dagBag.GetDagInstances() {
				fmt.Printf("start to send dag, %v", *dag)
				go func(dag *running.DAGInstance) {
					defer func() {
						if err := recover(); err != nil {
							//todo
							//e.dispatch.SendFail(dagInstance.TaskId, dagInstance.Id)
							return
						}
					}()
					//todo，校验是否dag已经分配
					workerInfo, err := e.workerManager.Select()
					if err != nil {
						//todo,重新发送dagbag
					}
					if workerInfo == nil {
						time.Sleep(time.Second)
						dagResend := dispatch.NewDagBag(dagBag.GetTask(), []*running.DAGInstance{dag})
						todoDags <- dagResend
						return
					}
					err = e.taskManager.SendTask(workerInfo.GetWorker().HostIp, workerInfo.GetWorker().Port, dag)
					if err != nil {
						//todo
						//e.dispatch.SendFail(dagInstance.TaskId, dagInstance.Id)
						return
					}
				}(dag)
			}
		}
	}
}
