package inspectors

import (
	"fmt"
	"iote"
	"iote/model/entities"
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

func (qi *QualityInspector) ReportService() *iote.EventInspectorService {
	return nil
}

func (qi *QualityInspector) Inspect(entity iote.EventEntity, errChan chan<- error) {
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
						errChan <- iote.EventAlert{
							Event: nil,
							Msg:   fmt.Sprintf("%s:%f第一次过高", k, wk),
						}
						errChan <- iote.EventAlert{
							Event: nil,
							Msg:   fmt.Sprintf("%s:%f第二次过高", k, v),
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
	} else {
		errChan <- fmt.Errorf("wrong entity type in QualityInspector!")
	}
}
