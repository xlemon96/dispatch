package task

import (
	"github.com/naivelife/dispatch/model/running"
)

type TaskManager interface {
	SendTask(ip, port string, dag *running.DAGInstance) error
}
