package bean

import (
	"time"
)

type Task struct {
	Id          int64     `gorm:"id"`
	TaskState   string    `gorm:"task_state"`
	TaskType    string    `gorm:"task_type"`
	CreatedTime time.Time `gorm:"created_time"`
	UpdateTime  time.Time `gorm:"update_time"`
}
