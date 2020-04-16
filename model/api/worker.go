package api

type HeartBeatRequest struct {
	HostIp            string
	Port              string
	RequestId         string `json:"request_id"`
	SchedulerEndpoint string `json:"scheduler_endpoint"`
}

type HeartBeatResponse struct {
	RunningTaskCount int
}

type SendTaskRequest struct {
	HostIp        string
	Port          string
	DagInstanceId int64  `json:"dag_instance_id"`
}

type SendTaskResponse struct {
}
