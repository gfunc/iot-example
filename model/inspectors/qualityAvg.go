package inspectors

import (
	"fmt"
	"iote"
	"iote/model/entities"
	"net/http"
	"strings"
	"sync"
	"time"
)

//1. 数据报表分析添加一个新的分析项：每日平均油品质量指标。输出如下：
//# 质量检测：
//2020-02-30,AB:xxx,AE:xxx,CE:xxx;
//2020-06-01,AB:xxx,AE:xxx,CE:xxx;

type QualityAvgInspector struct {
	sync.Mutex
	cache map[string]*qltAvg
}

type qltAvg struct {
	date       time.Time
	indexCache map[string]*qltIndexAvg
}

type qltIndexAvg struct {
	count int
	avg   float64
}

func (ti *QualityAvgInspector) ReportService() *iote.EventInspectorService {
	return &iote.EventInspectorService{
		Handler: func(w http.ResponseWriter, req *http.Request) {
			respStr := ""
			for t, qa := range ti.cache {
				indexStrs := make([]string, 0)
				for i, v := range qa.indexCache {
					indexStrs = append(indexStrs, fmt.Sprintf("%s:%f", i, v.avg))
				}
				indexStr := strings.Join(indexStrs, ",")
				respStr += fmt.Sprintf("设备: %s, %s, %s\n", t, qa.date.Format("2006-01-02"), indexStr)
			}
			fmt.Fprintf(w, respStr)
		},
		URI: "avg",
	}
}

func (ti *QualityAvgInspector) Inspect(entity iote.EventEntity, errChan chan<- error) {
	ti.Lock()
	defer ti.Unlock()
	if ti.cache == nil {
		ti.cache = make(map[string]*qltAvg, 0)
	}
	tp, ok := entity.(*entities.QualityEvent)
	if ok {
		tc, exist := ti.cache[entity.GetDeviceID()]
		if !exist {
			qa := &qltAvg{
				date:       tp.EventTime.Truncate(time.Hour * 24),
				indexCache: make(map[string]*qltIndexAvg, 0),
			}
			for i, v := range tp.Indexes {
				qa.indexCache[i] = &qltIndexAvg{
					count: 1,
					avg:   v,
				}
			}

			ti.cache[entity.GetDeviceID()] = qa
		} else {
			if tp.EventTime.Truncate(time.Hour * 24).After(tc.date) {
				//if tc.count > 0 {
				//	errChan <- iote.EventAlert{
				//		Event: nil,
				//		Msg:   fmt.Sprintf("设备 %s, %s,温度: %f", tp.DeviceID, tc.date.Format("2006-01-02"), tc.avgTemp),
				//	}
				//}
				tc.date = tp.EventTime
				for i, v := range tp.Indexes {
					tc.indexCache[i].count = 1
					tc.indexCache[i].avg = v
				}

				return
			}
			for i, v := range tp.Indexes {
				indexAvg := tc.indexCache[i]
				indexAvg.avg = ((indexAvg.avg * float64(indexAvg.count)) + v) / float64(indexAvg.count+1)
				indexAvg.count++
			}

		}

	} else {
		errChan <- fmt.Errorf("wrong entity type in QualityAvgInspector!")
	}
}
