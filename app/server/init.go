package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/naivelife/dispatch/constant"
	"github.com/naivelife/dispatch/engine"
	"github.com/naivelife/dispatch/util"
)

func initDB() *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		"root", "jiajianyun", "127.0.0.1", "3306", "dispatch")
	var db *gorm.DB
	var err error
	if db, err = gorm.Open("mysql", url); err != nil {
		panic(err)
	}
	//设置数据库参数
	//db.LogMode(true)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(10)
	return db
}

func initLog() *log.Logger {
	file, err := os.OpenFile("/Users/jiajianyun/go/src/github.com/navieboy/dispatch/log.txt", os.O_RDWR|os.O_APPEND, 777)
	if err != nil {
		panic(err)
	}
	logger := log.New(file, "[dispatch]", log.LstdFlags)
	return logger
}

func initWebServer(logger *log.Logger) (*http.Server, *util.Router) {
	router := util.NewRouter()
	router.Logger = logger
	router.RegisterFilters(router.PrepareFilter, router.InvokerFilter)
	server := &http.Server{
		Addr: "127.0.0.1:8899",
	}
	return server, router
}

func initHandler(s *Server, e *engine.Engine) {
	s.router.RegisterHandleFunc(constant.UpdateDAGInstance, e.GetDispatchCommunication().UpdateDAGInstance)
	s.router.RegisterHandleFunc(constant.DescribeDAGInstance, e.GetDispatchCommunication().DescribeDAGInstance)
	s.router.RegisterHandleFunc(constant.DescribeDAGInstances, e.GetDispatchCommunication().DescribeDAGInstances)
	s.router.RegisterHandleFunc(constant.DescribeDAGInstancesByTask, e.GetDispatchCommunication().DescribeDAGInstancesByTask)
	s.router.RegisterHandleFunc(constant.UpdateTask, e.GetDispatchCommunication().UpdateTask)
}
