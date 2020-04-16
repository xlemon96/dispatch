package task

import (
	"dispatch/model/running"
)

type TaskManager interface {
	SendTask(ip, port string, dag *running.DAGInstance) error
}
