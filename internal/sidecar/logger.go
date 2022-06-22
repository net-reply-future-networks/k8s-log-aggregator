package sidecar

import (
	"encoding/json"
	"fmt"
	"time"
)

type Logger struct {
	Container string
}

type Log struct {
	Process   *Pid   `json:"process,omitempty"`
	Container string `json:"container"`
	Log       string `json:"log"`
	Level     string `json:"level"`
	Time      int64  `json:"time"`
}

func (l *Logger) Log(log Log) {
	log.Container = l.Container
	log.Time = time.Now().Unix()
	b, err := json.Marshal(log)
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
