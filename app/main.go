package main

import (
	"runtime"
	"sync"

	"github.com/naivelife/dispatch/app/server"
	"github.com/naivelife/dispatch/util"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	s := server.NewServer()
	wg := &sync.WaitGroup{}
	tn := util.New(nil, s.Stop, func() {
		wg.Done()
	})
	if err := tn.Run(func() error {
		wg.Add(1)
		if err := s.Start(); err != nil {
			panic(err)
		}
		wg.Wait()
		return nil
	}); err != nil {
		panic(err)
	}
}