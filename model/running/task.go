package running

import (
	"time"
)

type Task struct {
	Id           int64
	Type         string
	State        string
	CreatedTime  time.Time
	UpdateTime   time.Time
}
