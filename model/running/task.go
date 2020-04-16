package running

import (
	"time"
)

type Task struct {
	Id           int64
	Type         string
	IncludeHosts []string
	ExcludeHosts []string
	Content      string
	RequestId    string
	Pin          string
	Source       string
	State        string
	CreatedTime  time.Time
	UpdateTime   time.Time
}
