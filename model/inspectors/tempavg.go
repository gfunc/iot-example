package inspectors

import (
	"fmt"
	ip "iot_practise"
	"iot_practise/model/entities"
	"sync"
	"time"
)

type TemperatureAvgInspector struct {
	sync.Mutex
	date    time.Time
	count   int
	avgTemp float64
}

func (ti *TemperatureAvgInspector) Inspect(entity ip.EventEntity, errChan chan<- error) {
	ti.Lock()
	defer ti.Unlock()

	tp, ok := entity.(*entities.TemperatureEvent)
	if ok {
		if tp.EventTime.Truncate(time.Hour * 24).After(ti.date) {
			if ti.count > 0 {
				errChan <- &ip.EventAlert{
					Event: nil,
					Msg:   fmt.Sprintf("温度: %s %f", ti.date.Format("2006-01-02"), ti.avgTemp),
				}
			}
			ti.count = 1
			ti.avgTemp = tp.Temperature
			return
		}

		ti.avgTemp = (ti.avgTemp * float64(ti.count)) / float64(ti.count+1)
		ti.count++
	}
}
