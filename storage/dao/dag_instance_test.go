package dao

import (
	"testing"
	"time"

	"github.com/navieboy/dispatch/model/bean"
)

var dagInstanceDao = &DagInstanceDaoImpl{}

func TestDagInstanceDaoImpl_Create(t *testing.T) {
	dagInstance1 := &bean.DagInstance{
		TaskId:      1,
		Depends:     "",
		HostIp:      "",
		Port:        "",
		DagType:     "",
		Output:      "",
		DagState:    "running",
		CreatedTime: time.Now(),
		UpdateTime:  time.Now(),
	}
	dagInstance2 := &bean.DagInstance{
		TaskId:      1,
		Depends:     "1",
		HostIp:      "",
		Port:        "",
		DagType:     "",
		Output:      "",
		DagState:    "running",
		CreatedTime: time.Now(),
		UpdateTime:  time.Now(),
	}
	dagInstance3 := &bean.DagInstance{
		TaskId:      1,
		Depends:     "2",
		HostIp:      "",
		Port:        "",
		DagType:     "",
		Output:      "",
		DagState:    "running",
		CreatedTime: time.Now(),
		UpdateTime:  time.Now(),
	}
	dagInstance4 := &bean.DagInstance{
		TaskId:      1,
		Depends:     "2",
		HostIp:      "",
		Port:        "",
		DagType:     "",
		Output:      "",
		DagState:    "running",
		CreatedTime: time.Now(),
		UpdateTime:  time.Now(),
	}
	if err := dagInstanceDao.Create(test_db, dagInstance1); err != nil {
		panic(err)
	}
	if err := dagInstanceDao.Create(test_db, dagInstance2); err != nil {
		panic(err)
	}
	if err := dagInstanceDao.Create(test_db, dagInstance3); err != nil {
		panic(err)
	}
	if err := dagInstanceDao.Create(test_db, dagInstance4); err != nil {
		panic(err)
	}
}