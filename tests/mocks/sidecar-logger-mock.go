package mocks

import (
	"fmt"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/sidecar"
)

type LogMock struct {
	LogOut        []sidecar.Log
	ErrorOut      []string
	ErrorPIDOut   []string
	InfoOut       []string
	InfoPIDOut    []string
	SidecarLogOut []string
}

func (lm *LogMock) Log(log sidecar.Log) {
	fmt.Println(log)
	lm.LogOut = append(lm.LogOut, log)
}
func (lm *LogMock) Error(log string) {
	fmt.Println(log)
	lm.ErrorOut = append(lm.ErrorOut, log)
}
func (lm *LogMock) ErrorPID(log string, pid sidecar.Pid) {
	fmt.Println(log)
	lm.ErrorPIDOut = append(lm.ErrorPIDOut, log)
}
func (lm *LogMock) Info(log string) {
	fmt.Println(log)
	lm.InfoOut = append(lm.InfoOut, log)
}
func (lm *LogMock) InfoPID(log string, pid sidecar.Pid) {
	fmt.Println(log)
	lm.InfoPIDOut = append(lm.InfoPIDOut, log)
}
func (lm *LogMock) SidecarLog(log string) {
	fmt.Println(log)
	lm.SidecarLogOut = append(lm.SidecarLogOut, log)
}
