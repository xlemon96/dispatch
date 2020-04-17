package running

import (
	"time"
)

type DAGInstance struct {
	Id           int64
	TaskId       int64
	HostIp       string
	Port         string
	Depends      []string
	State        string
	Type         string
	Output       string
	CreatedTime  time.Time
	UpdateTime   time.Time
}
