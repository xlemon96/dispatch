package bean

import (
	"time"
)

type Worker struct {
	Id          int64     `gorm:"id"`
	Name        string    `gorm:"name"`
	HostIp      string    `gorm:"host_ip"`
	Port        string    `gorm:"port"`
	WorkerState string    `gorm:"worker_state"`
	Disable     int32     `gorm:"disable"`
	CreatedTime time.Time `gorm:"created_time"`
	UpdateTime  time.Time `gorm:"update_time"`
}
