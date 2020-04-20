package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/navieboy/dispatch/engine"
	"github.com/navieboy/dispatch/storage"
	"github.com/navieboy/dispatch/util"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	e := initEngine()
	wg := &sync.WaitGroup{}
	tn := util.New(nil, e.Stop, func() {
		wg.Done()
	})
	if err := tn.Run(func() error {
		wg.Add(1)
		if err := e.Start(); err != nil {
			panic(err)
		}
		wg.Wait()
		return nil
	}); err != nil {
		panic(err)
	}
}

func initEngine() *engine.Engine {
	db := initDB()
	logger := initLog()
	webServer, router := initWebServer(logger)
	s := storage.NewStorageImpl(db)
	e := engine.NewEngine(s, logger, webServer, router)
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
	//db.LogMode(true)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(10)
	return db
}

func initLog() *log.Logger {
	file, err := os.OpenFile("/Users/jiajianyun/go/src/github.com/navieboy/dispatch/log.txt", os.O_RDWR | os.O_APPEND, 777)
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
