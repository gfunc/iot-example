package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"
)

const timeFormat = `2006-01-02 15:04:05`

func TestTemperature(t *testing.T) {
	wg := new(sync.WaitGroup)
	rand.Seed(65)
	for k := 1; k < 11; k++ {
		wg.Add(1)
		go func(group *sync.WaitGroup, index int) {
			defer group.Done()
			for i := 0; i < 1000; i++ {
				datediff := int(i / 100)
				data := fmt.Sprintf("T%d,%s,%f;", index, time.Now().AddDate(0, 0, datediff).Format(timeFormat), rand.Float64()+float64(rand.Intn(100)))
				post("tmp", data)
			}
		}(wg, k)
	}
	wg.Wait()
}

func TestQuality(t *testing.T) {
	wg := new(sync.WaitGroup)
	rand.Seed(75)
	for k := 1; k < 11; k++ {
		wg.Add(1)
		go func(group *sync.WaitGroup, index int) {
			defer group.Done()
			for i := 0; i < 1000; i++ {
				ab := rand.Float64() + float64(rand.Intn(100))
				ae := rand.Float64() + float64(rand.Intn(100))
				ce := rand.Float64() + float64(rand.Intn(100))
				data := fmt.Sprintf("Q%d,%s,AB:%f,AE:%f,CE:%f;", index, time.Now().Format(timeFormat), ab, ae, ce)
				post("qlt", data)
			}
		}(wg, k)
	}
	wg.Wait()
}

func TestExample(t *testing.T) {
	temp := []string{"T1,2020-01-30 19:00:01,22;",
		"T1,2020-01-30 19:00:00,25;",
		"T1,2020-01-30 19:00:02,28;"}

	qly := []string{"Q1,2020-01-30 19:30:10,AB:37.8,AE:100,CE:0.01;",
		"Q1,2020-01-30 19:30:20,AB:39.8,AE:100,CE:0.01;",
		"Q1,2020-01-30 19:30:25,AB:39.9,AE:100,CE:0.01;",
		"Q1,2020-01-30 19:30:32,AB:48.9,AE:101,CE:0.011;",
		"Q1,2020-01-30 19:30:40,AB:58.9,AE:103,CE:0.012;",
	}
	for _, t := range temp {
		post("tmp", t)
	}
	for _, t := range qly {
		post("qlt", t)
	}
}

func post(uri, data string) {
	rsp, err := http.Post("http://127.0.0.1:8088/"+uri, "text", bytes.NewBufferString(data))
	if err != nil {
		fmt.Println(err)
	} else {
		if rsp.StatusCode >= 300 {
			b, err1 := ioutil.ReadAll(rsp.Body)
			if err1 != nil {
				fmt.Println(rsp.StatusCode, err1)
			}
			fmt.Println(rsp.StatusCode, string(b))
		}

	}
}
