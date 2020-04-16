package dao

import (
	"dispatch/model/running"
)

type Storage interface {
	WorkerDao
	TaskDao
}

type WorkerDao interface {
	DescribeWorkers() ([]*running.Worker, error)           // 查询worker
	DescribeUpdatedWorkers(int) ([]*running.Worker, error) // 查询更新的worker
	UpdateWorker(name, state string) error                 // 更新worker状态
}

type TaskDao interface {
	DescribeTasks(state string) ([]*running.Task, error)                      // 查询不同状态的任务
	DescribeTask(taskId int64) (*running.Task, error)                         // 查询任务
	UpdateTaskState(id int64, state string) error                             // 更新任务状态
	DescribeDAGInstances(taskId int64) ([]*running.DAGInstance, error)        // 查询某任务所有DAG
	DescribeDAGInstance(dagId int64) (*running.DAGInstance, error)            // 查询单个DAG
	DescribeDAGInstancesByState(state string) ([]*running.DAGInstance, error) // 查询不同状态的任务
	UpdateDAGInstanceState(dagId int64, hostIp, port, state string) error     // 更新
	ClearDAGInstance(hostIp, port string)                                     // 更新hostIp的running的任务状态为rescheduling，重新开始分配
}
