package communication

import (
	"dispatch/model/api"
)

type WorkerClient interface {
	//Heartbeart模块获取worker心跳信息
	Heartbeat(request *api.HeartBeatRequest) (*api.HeartBeatResponse, error)
	//调度器下发任务到worker
	SendTask(request *api.SendTaskRequest) (*api.SendTaskResponse, error)
}
