package task

import (
	"github.com/navieboy/dispatch/model/running"
)

type TaskManager interface {
	SendTask(ip, port string, dag *running.DAGInstance) error
}
