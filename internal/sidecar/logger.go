package sidecar

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/utils"
)

type Logger struct {
	Pod        string
	NatsClient utils.NatsClientInterface
}

type LoggerInterface interface {
	Log(log Log)
	Error(log string)
	ErrorPID(log string, pid Pid)
	Info(log string)
	InfoPID(log string, pid Pid)
	SidecarLog(log string)
}

type Log struct {
	Process *Pid   `json:"process,omitempty"`
	Pod     string `json:"pod"`
	Log     string `json:"log"`
	Level   string `json:"level"`
	Time    int64  `json:"time"`
}

func (l *Logger) Log(log Log) {
	log.Pod = l.Pod
	log.Time = time.Now().Unix()
	b, err := json.Marshal(log)
	if err != nil {
		fmt.Println(err)
	}
	err = l.NatsClient.Publish(b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}

func (l *Logger) Error(log string) {
	logF := Log{
		Level: "error",
		Log:   log,
	}
	l.Log(logF)
}

func (l *Logger) ErrorPID(log string, pid Pid) {
	logF := Log{
		Process: &pid,
		Level:   "error",
		Log:     log,
	}
	l.Log(logF)
}

func (l *Logger) Info(log string) {
	logF := Log{
		Level: "info",
		Log:   log,
	}
	l.Log(logF)
}

func (l *Logger) InfoPID(log string, pid Pid) {
	logF := Log{
		Process: &pid,
		Level:   "info",
		Log:     log,
	}
	l.Log(logF)
}

func (l *Logger) SidecarLog(log string) {
	fmt.Println(log)
}
