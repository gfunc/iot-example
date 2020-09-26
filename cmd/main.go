package main

import (
	"context"
	"fmt"
	ip "iot_practise"
	"iot_practise/model/entities"
	"iot_practise/model/inspectors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const threadSize = 3

var monitors []*ip.EventMonitor

func init() {
	monitors = []*ip.EventMonitor{
		{
			Name:   "tmp",
			Entity: &entities.TemperatureEvent{},
			Inspectors: []ip.EventInspector{
				new(inspectors.TemperatureInspector),
				new(inspectors.TemperatureAvgInspector),
			},
		},

		{
			Name:   "qlt",
			Entity: &entities.QualityEvent{},
			Inspectors: []ip.EventInspector{
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
		go func(monitor *ip.EventMonitor, ctxt context.Context) {
			for index := range monitor.WorkerCreate {
				go monitor.GetWorker(monitor.Name)(ctxt, index, errChan)
			}
		}(m, ctx)

		handler := m.GetHandler()
		http.Handle("/"+m.Name, handler)
	}
	cn := ip.ConsoleNotifier{}
	go func() {
		for err := range errChan {
			ea, ok := err.(ip.EventAlert)
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
