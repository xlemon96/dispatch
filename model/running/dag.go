package running

import (
	"time"
)

type DAGInstance struct {
	Id           int64
	TaskId       int64
	HostIp       string
	Port         string
	IgnoreErr    bool
	Depends      []string
	State        string
	Type         string
	Output       string
	RetryCount   int
	Timeout      int
	CreatedTime  time.Time
	UpdateTime   time.Time
	ErrMsg       string
}
