package mocks

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/sidecar"
)

type StreamManagerInterface interface {
	OpenStream(pid sidecar.Pid)
	CloseStream(pid sidecar.Pid) bool
}

type SreamManagerMock struct {
	MockOpenStream  MockOpenStream
	MockCloseStream MockCloseStream

	T *testing.T
}

type MockOpenStream struct {
	InPid               sidecar.Pid
	ExpectedInvocations int
	ActualInvocations   int
}

type MockCloseStream struct {
	InPid               sidecar.Pid
	OutBool             bool
	ExpectedInvocations int
	ActualInvocations   int
}

func (sm *SreamManagerMock) OpenStream(pid sidecar.Pid) {
	sm.MockOpenStream.ActualInvocations++
	if !reflect.DeepEqual(sm.MockOpenStream.InPid, pid) {
		sm.T.Error(fmt.Sprintf("(OpenStream) unexpected input pid value, expected %v, got %v", sm.MockOpenStream.InPid, pid))
	}
}
func (sm *SreamManagerMock) CloseStream(pid sidecar.Pid) bool {
	sm.MockCloseStream.ActualInvocations++
	if !reflect.DeepEqual(sm.MockCloseStream.InPid, pid) {
		sm.T.Error(fmt.Sprintf("(CloseStream) unexpected input pid value, expected %v, got %v", sm.MockCloseStream.InPid, pid))
	}
	return sm.MockCloseStream.OutBool
}
