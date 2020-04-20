package engine

import (
	"log"

	"github.com/navieboy/dispatch/common"
	"github.com/navieboy/dispatch/engine/communication"
	"github.com/navieboy/dispatch/engine/dispatch"
	"github.com/navieboy/dispatch/engine/heartbeat"
	"github.com/navieboy/dispatch/engine/task"
	"github.com/navieboy/dispatch/engine/worker"
	"github.com/navieboy/dispatch/model/running"
	"github.com/navieboy/dispatch/storage"
)

type Engine struct {
	logger                *log.Logger
	dispatch              dispatch.Dispatch
	workerManager         worker.WorkerManager
	workerClient          communication.WorkerClient
	dispatchCommunication communication.DispatchCommunication
	taskManager           task.TaskManager
	heartbeat             heartbeat.Heartbeat
}

func NewEngine(dao storage.Storage, logger *log.Logger) *Engine {
	Engine := &Engine{
		dispatch:      dispatch.NewDispatch(dao, logger),
		workerManager: worker.NewWorkerManager(dao, logger),
		workerClient:  communication.NewWorkerClientImpl(common.NewDefaultClient(), logger),
	}
	Engine.dispatchCommunication = communication.NewDispatchCommunicationImpl(dao,
		Engine.dispatch, Engine.workerManager, logger)
	Engine.taskManager = task.NewTaskManager(dao, Engine.workerManager, Engine.workerClient, logger)
	Engine.heartbeat = heartbeat.NewHeartbeatImpl(Engine.workerClient, Engine.workerManager, logger)
	Engine.logger = logger
	return Engine
}

func (e *Engine) Start() error {
	var err error
	e.logger.Println("start engine......")
	if err = e.workerManager.Start(); err != nil {
		return err
	}
	if err = e.heartbeat.Start(); err != nil {
		return err
	}
	if err = e.dispatch.Start(); err != nil {
		return err
	}
	go e.startDispatch()
	e.logger.Println("start engine success......")
	return nil
}

func (e *Engine) Close() {
	//todo
}

func (e *Engine) GetDispatchCommunication() communication.DispatchCommunication {
	return e.dispatchCommunication
}

func (e *Engine) startDispatch() {
	todoDags := e.dispatch.GetTodoDags()
	for {
		select {
		case dagBag := <-todoDags:
			if len(dagBag.GetDagInstances()) == 0 {
				continue
			}
			for _, dag := range dagBag.GetDagInstances() {
				go func(dag *running.DAGInstance) {
					defer func() {
						if err := recover(); err != nil {
							//todo
							//e.dispatch.SendFail(dagInstance.TaskId, dagInstance.Id)
							return
						}
					}()
					//todo，校验是否dag已经分配
					workerInfo := e.workerManager.Select()
					if workerInfo == nil {
						dagResend := dispatch.NewDagBag(dagBag.GetTask(), []*running.DAGInstance{dag})
						todoDags <- dagResend
						return
					}
					err := e.taskManager.SendTask(workerInfo.GetWorker().HostIp, workerInfo.GetWorker().Port, dag)
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
