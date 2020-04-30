package dao

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/naivelife/dispatch/model/bean"
)

var test_db *gorm.DB
var err error
var taskDao = &TaskDaoImpl{}

func init()  {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		"root", "jiajianyun", "127.0.0.1", "3306", "dispatch")
	if test_db, err = gorm.Open("mysql", url); err != nil {
		panic(err)
	}
	//设置数据库参数
	test_db.LogMode(true)
	test_db.SingularTable(true)
	test_db.DB().SetMaxIdleConns(10)
	test_db.DB().SetMaxOpenConns(10)
}

func TestTaskDaoImpl_Create(t *testing.T) {
	task := &bean.Task{
		TaskState:   "running",
		TaskType:    "jcs",
		CreatedTime: time.Now(),
		UpdateTime:  time.Now(),
	}
	if err := taskDao.Create(test_db, task); err != nil {
		panic(err)
	}
}

func TestTaskDaoImpl_Query(t *testing.T) {
	id := int64(1)
	task, err := taskDao.Query(test_db, id)
	if err != nil {
		panic(err)
	}
	fmt.Println(task)
}

func TestTaskDaoImpl_QueryByState(t *testing.T) {
	state := "running"
	tasks, err := taskDao.QueryByState(test_db, state)
	if err != nil {
		panic(err)
	}
	fmt.Println(tasks)
}
