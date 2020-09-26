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
				data := fmt.Sprintf("T%d,%s,%f;", index, time.Now().AddDate(0, 0, int(i/100)).Format(timeFormat), rand.Float64()+float64(rand.Intn(100)))
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
