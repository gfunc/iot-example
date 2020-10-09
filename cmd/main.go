package main

import (
	"context"
	"fmt"
	"iote"
	"iote/model/entities"
	"iote/model/inspectors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const threadSize = 3

var monitors []*iote.EventMonitor

func init() {
	monitors = []*iote.EventMonitor{
		{
			Name:   "tmp",
			Entity: &entities.TemperatureEvent{},
			Inspectors: []iote.EventInspector{
				new(inspectors.TemperatureInspector),
				new(inspectors.TemperatureAvgInspector),
			},
		},

		{
			Name:   "qlt",
			Entity: &entities.QualityEvent{},
			Inspectors: []iote.EventInspector{
				new(inspectors.QualityInspector),
			},
		},
	}
}

func main() {
	cancelFuncs := make([]context.CancelFunc, 0)
	errChan := make(chan error, threadSize)
	killC := make(chan os.Signal)
	signal.Notify(killC, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-killC
		for _, cf := range cancelFuncs {
			cf()
		}
		close(errChan)
		os.Exit(1)
	}()

	for _, m := range monitors {
		m.CreateChan(threadSize)
		wk := m.GetWorker(m.Name)
		ctx, cancel := context.WithCancel(context.Background())
		cancelFuncs = append(cancelFuncs, cancel)
		for k := 0; k < threadSize; k++ {
			go wk(ctx, k, errChan)
		}
		// expire worker when time is up, to avoid too large file
		go func(monitor *iote.EventMonitor, ctxt context.Context) {
			for index := range monitor.WorkerCreate {
				go monitor.GetWorker(monitor.Name)(ctxt, index, errChan)
			}
		}(m, ctx)

		handler := m.GetHandler()
		http.Handle("/"+m.Name, handler)
	}
	cn := iote.ConsoleNotifier{}
	go func() {
		for err := range errChan {
			ea, ok := err.(iote.EventAlert)
			if ok {
				cn.Notify(ea)
			} else {
				fmt.Println("ERROR: " + err.Error())
			}
		}
	}()

	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		panic(err)
	}
}
