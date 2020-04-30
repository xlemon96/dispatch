package communication

import (
	"fmt"
	"log"

	"github.com/naivelife/dispatch/common"
	"github.com/naivelife/dispatch/constant"
	"github.com/naivelife/dispatch/model/api"
)

type workerClientImpl struct {
	logger *log.Logger
	client *common.HttpClient
}

func NewWorkerClientImpl(httpclient *common.HttpClient, logger *log.Logger) *workerClientImpl {
	return &workerClientImpl{
		logger: logger,
		client: httpclient,
	}
}

func (c *workerClientImpl) SendTask(request *api.SendTaskRequest) (*api.SendTaskResponse, error) {
	response := &api.SendTaskResponse{}
	url := generateUrl(request.HostIp, request.Port)
	err := c.client.CallHttpResponse(url, constant.SendTaskActionName, request, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *workerClientImpl) Heartbeat(request *api.HeartBeatRequest) (*api.HeartBeatResponse, error) {
	response := &api.HeartBeatResponse{}
	url := generateUrl(request.HostIp, request.Port)
	err := c.client.CallHttpResponse(url, constant.HeartBeatActionName, request, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func generateUrl(hostIp, port string) string {
	return fmt.Sprintf("http://%s:%s/%s", hostIp, port, "task-worker")
}
