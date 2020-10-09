package iote

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"time"
)

const persistDir = "log"

type Worker func(context.Context, int, chan<- error)

type EventMonitor struct {
	Name         string
	Entity       EventEntity
	entityChan   chan EventEntity
	Inspectors   []EventInspector
	WorkerCreate chan int
}

func (em *EventMonitor) CreateChan(size int) {
	em.entityChan = make(chan EventEntity, size)
	em.WorkerCreate = make(chan int, 1)
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
		em.entityChan <- e
		fmt.Fprintf(w, em.Name+" ok")
	}
}

func (em *EventMonitor) GetWorker(prefix string) Worker {
	etp := reflect.TypeOf(em.Entity)
	if etp.Kind() == reflect.Ptr {
		etp = etp.Elem()
	}

	if _, err := os.Stat(persistDir); os.IsNotExist(err) {
		os.Mkdir(persistDir, os.ModeDir)
	}
	return func(ctx context.Context, index int, errChan chan<- error) {
		pattern := fmt.Sprintf("%s_%d_%s", prefix, index, etp.Name())
		f, err := ioutil.TempFile(persistDir, pattern)
		if err != nil {
			errChan <- err
			return
		}

		timer := time.NewTimer(5*time.Hour + time.Duration(index)*time.Minute)

		defer f.Sync()
		defer f.Close()
		defer timer.Stop()
		for e := range em.entityChan {
			select {
			case <-ctx.Done():
				em.entityChan <- e
				return
			case <-timer.C:
				em.WorkerCreate <- index
				em.entityChan <- e
				return
			default:
				go func(entity EventEntity) {
					for _, i := range em.Inspectors {
						i.Inspect(entity, errChan)
					}
				}(e)

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
