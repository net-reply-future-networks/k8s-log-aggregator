package sidecar

import (
	"encoding/json"
	"fmt"
)

type Logger struct {
	Container string
}

type Log struct {
	Process   *Pid   `json:"process,omitempty"`
	Container string `json:"container"`
	Log       string `json:"log"`
	Level     string `json:"level"`
}

func (l *Logger) Log(log Log) {
	log.Container = l.Container
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
