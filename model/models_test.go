package model

import (
	"encoding/json"
	"fmt"
	"iote/model/entities"
	"testing"
	"time"
)

func TestEvent_ParseBinary(t *testing.T) {
	raw := []string{
		`T1,2020-01-30 19:00:00,25.8;`,
		`T1,2020-01-30 19:00:01,22;`,
		`T1,2020-01-30 19:00:02,28;`,
		`Q1,2020-01-30 19:30:10,AB:37.8,AE:100,CE:0.01;`,
		`Q1,2020-01-30 19:30:20,AB:39.8,AE:100,CE:0.01;`,
		`Q1,2020-01-30 19:30:25,AB:39.9,AE:100,CE:0.01;`,
	}

	for _, s := range raw {
		e := entities.Event{}
		err := e.ParseBinary([]byte(s))
		if err != nil {
			panic(err)
		}
		printOut(e)
	}
}

func TestTemperatureEvent_ParseBinary(t *testing.T) {
	raw := []string{
		`T1,2020-01-30 19:00:00,25.8;`,
		`T1,2020-01-30 19:00:01,22;`,
		`T1,2020-01-30 19:00:02,28;`,
	}

	for _, s := range raw {
		e := entities.TemperatureEvent{}
		err := e.ParseBinary([]byte(s))
		if err != nil {
			panic(err)
		}
		printOut(e)
	}
}

func TestQualityEvent_ParseBinary(t *testing.T) {
	raw := []string{
		`Q1,2020-01-30 19:30:10,AB:37.8,AE:100,CE:0.01;`,
		`Q1,2020-01-30 19:30:20,AB:39.8,AE:100,CE:0.01;`,
		`Q1,2020-01-30 19:30:25,AB:39.9,AE:100,CE:0.01;`,
	}

	for _, s := range raw {
		e := entities.QualityEvent{}
		err := e.ParseBinary([]byte(s))
		if err != nil {
			panic(err)
		}
		printOut(e)
	}
}

func printOut(i interface{}) {
	b, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func TestTimer(t *testing.T) {
	timer := time.NewTimer(2 * time.Second)
	tt := <-timer.C
	fmt.Printf("time is up %s", tt.Format("2006-01-02T15:04:05"))
}
