package inspectors

import (
	ip "iot_practise"
	"iot_practise/model/entities"
	"sync"
	"time"
)

type TemperatureInspector struct {
	sync.Mutex
	cache map[string]tempStat
}

type tempStat struct {
	date time.Time
	max  entities.TemperatureEvent
	min  entities.TemperatureEvent
}

func (ti *TemperatureInspector) Inspect(entity ip.EventEntity, errChan chan<- error) {
	ti.Lock()
	defer ti.Unlock()
	if ti.cache == nil {
		ti.cache = make(map[string]tempStat, 0)
	}

	tp, ok := entity.(*entities.TemperatureEvent)

	if ok {
		tc, exist := ti.cache[entity.GetDeviceID()]
		if !exist {
			ti.cache[entity.GetDeviceID()] = tempStat{
				date: tp.EventTime.Truncate(time.Hour * 24),
				max:  *tp,
				min:  *tp,
			}
		} else {
			if tp.EventTime.Truncate(time.Hour * 24).After(tc.date) {
				tc.max = *tp
				tc.min = *tp
				return
			}
			msg := ""
			if tp.Temperature > tc.max.Temperature {
				tc.max = *tp
				if tp.Temperature-tc.min.Temperature >= 5 {
					msg = "温度过高"
				}
			}
			if tp.Temperature < tc.min.Temperature {
				tc.min = *tp
				if tc.max.Temperature-tp.Temperature >= 5 {
					msg = "温度过低"
				}
			}
			if msg != "" {
				errChan <- &ip.EventAlert{
					Event: entity,
					Msg:   msg,
				}
			}
		}
	}
}
