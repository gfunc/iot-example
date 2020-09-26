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

var monitors []ip.EventMonitor

func init() {
	monitors = []ip.EventMonitor{
		{
			Name:       "tmp",
			Entity:     &entities.TemperatureEvent{},
			EntityChan: make(chan ip.EventEntity, threadSize),
			Inspectors: []ip.EventInspector{
				new(inspectors.TemperatureInspector),
				new(inspectors.TemperatureAvgInspector),
			},
		},

		{
			Name:       "qlt",
			Entity:     &entities.QualityEvent{},
			EntityChan: make(chan ip.EventEntity, threadSize),
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
		http.HandleFunc("/"+m.Name, m.GetHandler())
		wk := m.GetWorker(m.Name)
		ctx, cancel := context.WithCancel(context.Background())
		cancelFuncs = append(cancelFuncs, cancel)
		for i := 0; i < threadSize; i++ {
			go wk(ctx, errChan)
		}
	}
	go func() {
		for err := range errChan {
			ea, ok := err.(*ip.EventAlert)
			if ok {
				fmt.Println("WARNING: " + ea.Error())
			} else {
				fmt.Println("ERROR: " + ea.Error())
			}
		}
	}()

	err := http.ListenAndServe(":8088", nil)
	if err != nil {
		panic(err)
	}
}
