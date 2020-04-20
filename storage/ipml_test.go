package storage

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/navieboy/dispatch/storage/dao"
)

var db = createDB()
var storage = &storageImpl{
	TaskDaoImpl:        dao.TaskDaoImpl{},
	DagInstanceDaoImpl: dao.DagInstanceDaoImpl{},
	WorkerDaoImpl:      dao.WorkerDaoImpl{},
	db:                 db,
}

func createDB() *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		"root", "jiajianyun", "127.0.0.1", "3306", "dispatch")
	var test_db *gorm.DB
	var err error
	if test_db, err = gorm.Open("mysql", url); err != nil {
		panic(err)
	}
	//设置数据库参数
	test_db.LogMode(true)
	test_db.SingularTable(true)
	test_db.DB().SetMaxIdleConns(10)
	test_db.DB().SetMaxOpenConns(10)
	return test_db
}

func TestStorageImpl_DescribeTask(t *testing.T) {
	id := int64(1)
	task, err := storage.DescribeTask(id)
	if err != nil {
		panic(nil)
	}
	fmt.Println(task)
}

func TestStorageImpl_DescribeTasks(t *testing.T) {
	state := "running"
	tasks, err := storage.DescribeTasks(state)
	if err != nil {
		panic(nil)
	}
	fmt.Println(tasks)
}

func TestStorageImpl_DescribeWorkers(t *testing.T) {
	workers, err := storage.DescribeWorkers()
	if err != nil {
		panic(err)
	}
	fmt.Println(workers)
}