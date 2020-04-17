package api

type UpdateDAGInstance struct {
	RequestId string
	TaskId    int64
	DAGId     int64 // dag instance id
	HostIp    string
	Port      string
	State     string
	OutPut    string
	RollBack  bool
	ErrMsg    string
}

type DagDependence struct {
	DagInstanceId int64
	DagId         int64
	Type          string
	Output        string
}

type DagInstance struct {
	TaskId        int64
	DagInstanceId int64
	DagId         int64
	ReferId       string // 引用id，可以为空
	HostIp        string
	Port          string
	Type          string
	State         string //refer to constant.DagTaskState...
	Output        string
	Dependence    []*DagDependence
	RequestId     string
	Content       string
	ErrMsg        string
}

type UpdateTask struct {
	TaskId    int64
	State     string
}

type Test struct {
	Test int `json:"test"`
}
