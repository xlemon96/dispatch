package util

import (
	"errors"
	"sync"
)

type LaunchFunc func() error

type FinishFunc func()

type Server struct {
	stop       chan struct{}
	isStart    bool
	mut        *sync.Mutex
	launchFunc LaunchFunc
	finishFunc FinishFunc
}

func NewServer(launchFunc LaunchFunc, finishFunc FinishFunc) *Server {
	return &Server{
		stop:       make(chan struct{}),
		isStart:    false,
		mut:        &sync.Mutex{},
		launchFunc: launchFunc,
		finishFunc: finishFunc,
	}
}

func (s *Server) Start() error {
	s.Lock()
	if s.isStart {
		s.Unlock()
		return errors.New("service is started")
	}
	s.isStart = true
	s.Unlock()
	if s.launchFunc != nil {
		return s.launchFunc()
	}
	return nil
}

func (s *Server) Stop() {
	s.Lock()
	if !s.isStart {
		s.Unlock()
		return
	}
	s.isStart = false
	close(s.stop)
	s.Unlock()
	if s.finishFunc != nil {
		s.finishFunc()
	}
}

func (s *Server) IsStart() bool {
	return s.isStart
}

func (s *Server) Lock() {
	s.mut.Lock()
}

func (s *Server) Unlock() {
	s.mut.Unlock()
}

func (s *Server) Close() chan struct{} {
	return s.stop
}
