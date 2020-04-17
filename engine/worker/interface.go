package worker

type WorkerManager interface {
	Start() error
	WorkerLoadIncrease(hostIp, port string)
	Select() *WorkerInfo
	ListAllWorkers() ([]*WorkerInfo, error)
	UpdateWorkerState(name, state string) error
}
