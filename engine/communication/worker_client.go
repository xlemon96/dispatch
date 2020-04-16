package communication

import (
	"fmt"

	"dispatch/common"
	"dispatch/constant"
	"dispatch/model/api"
)

type workerClientImpl struct {
	client *common.HttpClient
}

func NewWorkerClientImpl(httpclient *common.HttpClient) *workerClientImpl {
	return &workerClientImpl{client:httpclient}
}

//调度器下发任务到worker
func (c *workerClientImpl) SendTask(request *api.SendTaskRequest) (*api.SendTaskResponse, error) {
	response := &api.SendTaskResponse{}
	url := generateUrl(request.HostIp, request.Port)
	err := c.client.CallHttpResponse(url, constant.SendTaskActionName, request, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *workerClientImpl) Heartbeat(request *api.HeartBeatRequest) (*api.HeartBeatResponse, error)  {
	return nil, nil
}

func generateUrl(hostIp, port string) string {
	return fmt.Sprintf("http://%s:%s/%s", hostIp, port, "task-worker")
}