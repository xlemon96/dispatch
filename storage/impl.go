package storage

import (
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/naivelife/dispatch/model/bean"
	"github.com/naivelife/dispatch/model/running"
	"github.com/naivelife/dispatch/storage/dao"
)

type storageImpl struct {
	dao.TaskDaoImpl
	dao.DagInstanceDaoImpl
	dao.WorkerDaoImpl
	db *gorm.DB
}

func NewStorageImpl(db *gorm.DB) *storageImpl {
	return &storageImpl{
		TaskDaoImpl:        dao.TaskDaoImpl{},
		DagInstanceDaoImpl: dao.DagInstanceDaoImpl{},
		WorkerDaoImpl:      dao.WorkerDaoImpl{},
		db:                 db,
	}
}

func (s *storageImpl) DescribeTask(taskId int64) (*running.Task, error) {
	task, err := s.TaskDaoImpl.Query(s.db, taskId)
	if err != nil {
		return nil, err
	}
	rtask := &running.Task{
		Id:           task.Id,
		Type:         task.TaskType,
		State:        task.TaskState,
		CreatedTime:  task.CreatedTime,
		UpdateTime:   task.UpdateTime,
	}
	return rtask, nil
}

func (s *storageImpl) DescribeTasks(state string) ([]*running.Task, error) {
	var rtasks = make([]*running.Task, 0)
	tasks, err := s.TaskDaoImpl.QueryByState(s.db, state)
	if err != nil {
		return nil, err
	}
	for _, task := range tasks	 {
		rtasks = append(rtasks, &running.Task{
			Id:           task.Id,
			Type:         task.TaskType,
			State:        task.TaskState,
			CreatedTime:  task.CreatedTime,
			UpdateTime:   task.UpdateTime,
		})
	}
	return rtasks, nil
}

func (s *storageImpl) UpdateTaskState(id int64, state string) error {
	return nil
}

func (s *storageImpl) DescribeDAGInstance(dagId int64) (*running.DAGInstance, error) {
	dagInstance, err := s.DagInstanceDaoImpl.Query(s.db, dagId)
	if err != nil {
		return nil, err
	}
	rdagInstance := &running.DAGInstance{
		Id:          dagInstance.Id,
		TaskId:      dagInstance.TaskId,
		HostIp:      dagInstance.HostIp,
		Port:        dagInstance.Port,
		Depends:     strConvert2Slice(dagInstance.Depends),
		State:       dagInstance.DagState,
		Type:        dagInstance.DagType,
		Output:      dagInstance.Output,
		CreatedTime: dagInstance.CreatedTime,
		UpdateTime:  dagInstance.UpdateTime,
	}
	return rdagInstance, nil
}

func (s *storageImpl) DescribeDAGInstancesByState(state string) ([]*running.DAGInstance, error) {
	dagInstances, err := s.DagInstanceDaoImpl.QueryByState(s.db, state)
	if err != nil {
		return nil, err
	}
	rdagInstances := make([]*running.DAGInstance, 0)
	for _, dagInstance := range dagInstances {
		rdagInstances = append(rdagInstances, &running.DAGInstance{
			Id:          dagInstance.Id,
			TaskId:      dagInstance.TaskId,
			HostIp:      dagInstance.HostIp,
			Port:        dagInstance.Port,
			Depends:     strConvert2Slice(dagInstance.Depends),
			State:       dagInstance.DagState,
			Type:        dagInstance.DagType,
			Output:      dagInstance.Output,
			CreatedTime: dagInstance.CreatedTime,
			UpdateTime:  dagInstance.UpdateTime,
		})
	}
	return rdagInstances, nil
}

func (s *storageImpl) DescribeDAGInstances(taskId int64) ([]*running.DAGInstance, error) {
	dagInstances, err := s.DagInstanceDaoImpl.QueryByTaskID(s.db, taskId)
	if err != nil {
		return nil, err
	}
	rdagInstances := make([]*running.DAGInstance, 0)
	for _, dagInstance := range dagInstances {
		rdagInstances = append(rdagInstances, &running.DAGInstance{
			Id:          dagInstance.Id,
			TaskId:      dagInstance.TaskId,
			HostIp:      dagInstance.HostIp,
			Port:        dagInstance.Port,
			Depends:     strConvert2Slice(dagInstance.Depends),
			State:       dagInstance.DagState,
			Type:        dagInstance.DagType,
			Output:      dagInstance.Output,
			CreatedTime: dagInstance.CreatedTime,
			UpdateTime:  dagInstance.UpdateTime,
		})
	}
	return rdagInstances, nil
}

func (s *storageImpl) UpdateDAGInstanceState(dagId int64, hostIp, port, state string) error {
	return nil
}

func (s *storageImpl) ClearDAGInstance(hostIp, port string) {

}

func (s *storageImpl) DescribeWorkers() ([]*running.Worker, error) {
	rWorkers := make([]*running.Worker, 0)
	workers, err := s.WorkerDaoImpl.List(s.db, &bean.Worker{})
	if err != nil {
		return nil, err
	}
	for _, worker := range workers {
		rWorkers = append(rWorkers, &running.Worker{
			Id:          worker.Id,
			HostIp:      worker.HostIp,
			Port:        worker.Port,
			WorkerState: worker.WorkerState,
			Name:        worker.Name,
			Disable:     worker.Disable > 0,
		})
	}
	return rWorkers, nil
}

func (s *storageImpl) DescribeUpdatedWorkers(int) ([]*running.Worker, error) {
	return nil, nil
}

func (s *storageImpl) UpdateWorker(name, state string) error {
	return nil
}

func strConvert2Slice(input string) []string {
	if input == "" {
		return nil
	}
	if !strings.Contains(input, ",") {
		return []string{input}
	}
	return strings.Split(input, ",")
}
