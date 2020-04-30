package communication

import (
	"fmt"
	"log"

	"github.com/naivelife/dispatch/engine/dispatch"
	"github.com/naivelife/dispatch/engine/worker"
	"github.com/naivelife/dispatch/model/api"
	"github.com/naivelife/dispatch/storage"
	"github.com/naivelife/dispatch/util"
)

type dispatchCommunicationImpl struct {
	logger        *log.Logger
	storage       storage.Storage
	dispatch      dispatch.Dispatch
	workerManager worker.WorkerManager
}

func NewDispatchCommunicationImpl(storage storage.Storage, dispatch dispatch.Dispatch,
	workerManager worker.WorkerManager, logger *log.Logger) *dispatchCommunicationImpl {
	return &dispatchCommunicationImpl{
		logger:        logger,
		storage:       storage,
		dispatch:      dispatch,
		workerManager: workerManager,
	}
}

func (c *dispatchCommunicationImpl) UpdateDAGInstance(request *util.Request, response *util.Response) (err error) {
	req := request.BusinessRequest.(*api.Test)
	fmt.Println(req.Test)
	return nil
}

// 根据hostIp，查询
func (c *dispatchCommunicationImpl) DescribeDAGInstances(request *util.Request, response *util.Response) (err error) {
	return nil
}

func (c *dispatchCommunicationImpl) DescribeDAGInstancesByTask(request *util.Request, response *util.Response) (err error) {
	return nil
}

func (c *dispatchCommunicationImpl) DescribeDAGInstance(request *util.Request, response *util.Response) (err error) {
	return nil
}

// update task state
func (c *dispatchCommunicationImpl) UpdateTask(request *util.Request, response *util.Response) (err error) {
	return nil
}
