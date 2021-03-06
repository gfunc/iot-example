package iote

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EventEntity interface {
	GetDeviceID() string
	ParseBinary([]byte) error
}

type EventAlert struct {
	Event EventEntity
	Msg   string
}

func (ea EventAlert) Error() string {
	if ea.Event == nil {
		return ea.Msg
	}
	b, err := json.Marshal(ea.Event)
	if err != nil {
		return "error marshal json for event alert" + err.Error()
	}
	return string(b) + " " + ea.Msg
}

type EventInspector interface {
	Inspect(entity EventEntity, errChan chan<- error)
	ReportService() *EventInspectorService
}

type EventInspectorService struct {
	Handler http.HandlerFunc
	URI     string
}

type EventNotifier interface {
	Notify(EventAlert)
}

type ConsoleNotifier struct{}

func (ConsoleNotifier) Notify(ea EventAlert) {
	fmt.Println("WARNING: " + ea.Error())
}
