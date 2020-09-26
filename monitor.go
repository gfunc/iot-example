package iot_practise

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

type Worker func(context.Context, chan<- error)

type EventMonitor struct {
	Name       string
	Entity     EventEntity
	EntityChan chan EventEntity
	Inspectors []EventInspector
}

func (em *EventMonitor) GetHandler() http.HandlerFunc {
	etp := reflect.TypeOf(em.Entity)
	if etp.Kind() == reflect.Ptr {
		etp = etp.Elem()
	}
	return func(w http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}
		e := reflect.New(etp).Interface().(EventEntity)
		err = e.ParseBinary(b)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}
		em.EntityChan <- e
		fmt.Fprintf(w,"ok")
	}
}

func (em *EventMonitor) GetWorker(prefix string) Worker {
	etp := reflect.TypeOf(em.Entity)
	if etp.Kind() == reflect.Ptr {
		etp = etp.Elem()
	}
	return func(ctx context.Context, errChan chan<- error) {
		f, err := ioutil.TempFile(".", prefix+"_"+etp.Name())
		if err != nil {
			errChan <- err
			return
		}
		defer f.Sync()
		defer f.Close()
		select {
		case <-ctx.Done():
			return

		default:
			for e := range em.EntityChan {
				for _, i := range em.Inspectors {
					i.Inspect(e, errChan)
				}
				b, err := json.Marshal(e)
				if err != nil {
					errChan <- err
				}
				b = append(b, []byte("\n")...)
				_, err = f.Write(b)
				if err != nil {
					errChan <- err
				}
			}
		}
	}
}
