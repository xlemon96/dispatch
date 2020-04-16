package running

type Worker struct {
	Id          int64
	HostIp      string
	Port        string
	WorkerState string
	WorkerType  string
	Name        string
	Disable     bool
}
