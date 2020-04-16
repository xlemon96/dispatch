package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"dispatch/engine"
	"dispatch/storage"
)

func main() {
	waitStop(initEnv())
}

func initEnv() *engine.Engine {
	db := initDB()
	s := storage.NewStorageImpl(db)
	e := engine.NewEngine(s)
	err := e.Start()
	if err != nil {
		panic(err)
	}
	return e
}

func initDB() *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		"root", "jiajianyun", "127.0.0.1", "3306", "dispatch")
	var db *gorm.DB
	var err error
	if db, err = gorm.Open("mysql", url); err != nil {
		panic(err)
	}
	//设置数据库参数
	db.LogMode(true)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(10)
	return db
}

func waitStop(e *engine.Engine) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-signalChan
	fmt.Println("byebye")
}