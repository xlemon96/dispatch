package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/navieboy/dispatch/engine"
	"github.com/navieboy/dispatch/storage"
	"github.com/navieboy/dispatch/util"
)

type Server struct {
	*util.Server
	logger *log.Logger
	server *http.Server
	router *util.Router
	engine *engine.Engine
}

func NewServer() *Server {
	logger := initLog()
	webServer, router := initWebServer(logger)
	db := initDB()
	dao := storage.NewStorageImpl(db)
	e := engine.NewEngine(dao, logger)
	s := &Server{
		logger: logger,
		server: webServer,
		router: router,
		engine: e,
	}
	s.Server = util.NewServer(s.doStart, s.doClose)
	initHandler(s, e)
	return s
}

func (s *Server) doStart() error {
	var err error
	if err = s.engine.Start(); err != nil {
		return err
	}
	go func() {
		http.Handle(fmt.Sprintf("/%s", "task"), s.router)
		err = s.server.ListenAndServe()
	}()
	time.Sleep(time.Second)
	if err != nil {
		s.logger.Println("start fail......")
		return err
	}
	return nil
}

func (s *Server) doClose() {
	s.logger.Println("stop success......")
}