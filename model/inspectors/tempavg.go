package inspectors

import (
	"fmt"
	"iote"
	"iote/model/entities"
	"sync"
	"time"
)

type TemperatureAvgInspector struct {
	sync.Mutex
	cache map[string]*tempAvg
}

type tempAvg struct {
	date    time.Time
	count   int
	avgTemp float64
}

func (ti *TemperatureAvgInspector) Inspect(entity iote.EventEntity, errChan chan<- error) {
	ti.Lock()
	defer ti.Unlock()
	if ti.cache == nil {
		ti.cache = make(map[string]*tempAvg, 0)
	}
	tp, ok := entity.(*entities.TemperatureEvent)
	if ok {
		tc, exist := ti.cache[entity.GetDeviceID()]
		if !exist {
			ti.cache[entity.GetDeviceID()] = &tempAvg{
				date:    tp.EventTime.Truncate(time.Hour * 24),
				count:   1,
				avgTemp: tp.Temperature,
			}
		} else {
			if tp.EventTime.Truncate(time.Hour * 24).After(tc.date) {
				if tc.count > 0 {
					errChan <- iote.EventAlert{
						Event: nil,
						Msg:   fmt.Sprintf("设备 %s, %s,温度: %f", tp.DeviceID, tc.date.Format("2006-01-02"), tc.avgTemp),
					}
				}
				tc.date = tp.EventTime
				tc.count = 1
				tc.avgTemp = tp.Temperature
				return
			}

			tc.avgTemp = ((tc.avgTemp * float64(tc.count)) + tp.Temperature) / float64(tc.count+1)
			tc.count++
		}

	} else {
		errChan <- fmt.Errorf("wrong entity type in TemperatureAvgInspector!")
	}
}
