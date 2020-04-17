package engine

import (
	"fmt"
	"log"
	"net/http"

	"dispatch/common"
	"dispatch/constant"
	"dispatch/engine/communication"
	"dispatch/engine/dispatch"
	"dispatch/engine/heartbeat"
	"dispatch/engine/task"
	"dispatch/engine/worker"
	"dispatch/model/running"
	"dispatch/storage"
	"dispatch/util"
)

type Engine struct {
	*util.Server
	logger                *log.Logger
	server                *http.Server
	router                *util.Router
	dispatch              dispatch.Dispatch
	workerManager         worker.WorkerManager
	workerClient          communication.WorkerClient
	dispatchCommunication communication.DispatchCommunication
	taskManager           task.TaskManager
	heartbeat             heartbeat.Heartbeat
}

func NewEngine(dao storage.Storage, logger *log.Logger, server *http.Server, router *util.Router) *Engine {
	Engine := &Engine{
		dispatch:      dispatch.NewDispatch(dao, logger),
		workerManager: worker.NewWorkerManager(dao, logger),
		workerClient:  communication.NewWorkerClientImpl(common.NewDefaultClient(), logger),
	}
	Engine.dispatchCommunication = communication.NewDispatchCommunicationImpl(dao,
		Engine.dispatch, Engine.workerManager, logger)
	Engine.taskManager = task.NewTaskManager(dao, Engine.workerManager, Engine.workerClient, logger)
	Engine.heartbeat = heartbeat.NewHeartbeatImpl(Engine.workerClient, Engine.workerManager, logger)
	Engine.server = server
	Engine.router = router
	Engine.logger = logger
	Engine.Server = util.NewServer(Engine.doStart, Engine.doClose)
	return Engine
}

func (e *Engine) doStart() error {
	e.logger.Println("start engine......")
	e.initHandler()
	if err := e.workerManager.Start(); err != nil {
		return err
	}
	if err := e.heartbeat.Start(); err != nil {
		return err
	}
	if err := e.dispatch.Start(); err != nil {
		return err
	}
	go e.startDispatch()

	http.Handle(fmt.Sprintf("/%s", "task"), e.router)
	if err := e.server.ListenAndServe(); err != nil {
		return err
	}
	e.logger.Println("start engine success......")
	return nil
}

func (e *Engine) doClose() {
	//todo
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

func (e *Engine) initHandler() {
	e.router.RegisterHandleFunc(constant.UpdateDAGInstance, e.dispatchCommunication.UpdateDAGInstance)
	e.router.RegisterHandleFunc(constant.DescribeDAGInstance, e.dispatchCommunication.DescribeDAGInstance)
	e.router.RegisterHandleFunc(constant.DescribeDAGInstances, e.dispatchCommunication.DescribeDAGInstances)
	e.router.RegisterHandleFunc(constant.DescribeDAGInstancesByTask, e.dispatchCommunication.DescribeDAGInstancesByTask)
	e.router.RegisterHandleFunc(constant.UpdateTask, e.dispatchCommunication.UpdateTask)
}
