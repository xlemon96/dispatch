package bean

import (
	"time"
)

type DagInstance struct {
	Id           int64     `gorm:"id"`
	TaskId       int64     `gorm:"task_id"`
	Depends      string    `gorm:"depends"`
	HostIp       string    `gorm:"host_ip"`
	Port         string    `gorm:"port"`
	DagType      string    `gorm:"dag_type"`
	Output       string    `gorm:"output"`
	DagState     string    `gorm:"dag_state"`
	CreatedTime  time.Time `gorm:"created_time"`
	UpdateTime   time.Time `gorm:"update_time"`
}
