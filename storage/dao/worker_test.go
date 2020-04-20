package dao

import (
	"fmt"
	"testing"
	"time"

	"github.com/navieboy/dispatch/model/bean"
)

var workerDao = &WorkerDaoImpl{}

func TestWorkerDaoImpl_Create(t *testing.T) {
	worker := &bean.Worker{
		Name:        "test_worker1",
		HostIp:      "10.10.10.1",
		Port:        "8888",
		WorkerState: "running",
		Disable:     1,
		CreatedTime: time.Now(),
		UpdateTime:  time.Now(),
	}
	if err := workerDao.Create(test_db, worker); err != nil {
		panic(err)
	}
}

func TestWorkerDaoImpl_Query(t *testing.T) {
	id := int64(1)
	worker, err := taskDao.Query(test_db, id)
	if err != nil {
		panic(err)
	}
	fmt.Println(worker)
}

func TestWorkerDaoImpl_List(t *testing.T) {
	//param若所有参数为空值，则返回所有记录
	//若param为默认值，则此参数不参与sql拼接
	param := &bean.Worker{
		WorkerState: "running",
		Disable:     1,
	}
	workers, err := workerDao.List(test_db, param)
	if err != nil {
		panic(err)
	}
	fmt.Println(workers)
}
