package worker

type WorkerManager interface {
	Start() error
	WorkerLoadIncrease(hostIp, port string)
	Select() (*WorkerInfo, error)
}
