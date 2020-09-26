package inspectors

import (
	"fmt"
	ip "iot_practise"
	"iot_practise/model/entities"
	"sync"
)

type QualityInspector struct {
	sync.Mutex
	cache map[string]indexStat
}

type indexStat struct {
	last map[string]float64
	warn map[string]float64
}

func (qi *QualityInspector) Inspect(entity ip.EventEntity, errChan chan<- error) {
	qi.Lock()
	defer qi.Unlock()
	if qi.cache == nil {
		qi.cache = make(map[string]indexStat, 0)
	}
	tp, ok := entity.(*entities.QualityEvent)
	if ok {
		qe, exist := qi.cache[entity.GetDeviceID()]
		if !exist {
			qi.cache[entity.GetDeviceID()] = indexStat{
				last: tp.Indexes,
				warn: make(map[string]float64),
			}
		} else {
			for k, v := range tp.Indexes {
				if wk, ok := qe.warn[k]; ok {
					if v >= wk*1.1 {
						errChan <- &ip.EventAlert{
							Event: entity,
							Msg:   fmt.Sprintf("%s过高", k),
						}
					} else {
						delete(qe.warn, k)
					}
				} else {
					if v >= qe.last[k]*1.1 {
						qe.warn[k] = v
					}
				}
				qe.last[k] = v
			}
		}
	}
}
