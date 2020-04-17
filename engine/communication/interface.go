package communication

import (
	"dispatch/model/api"
	"dispatch/util"
)

type DispatchCommunication interface {
	UpdateDAGInstance(request *util.Request, response *util.Response) (err error)
	DescribeDAGInstances(request *util.Request, response *util.Response) (err error)
	DescribeDAGInstancesByTask(request *util.Request, response *util.Response) (err error)
	DescribeDAGInstance(request *util.Request, response *util.Response) (err error)
	UpdateTask(request *util.Request, response *util.Response) (err error)
}

type WorkerClient interface {
	Heartbeat(request *api.HeartBeatRequest) (*api.HeartBeatResponse, error)
	SendTask(request *api.SendTaskRequest) (*api.SendTaskResponse, error)
}
