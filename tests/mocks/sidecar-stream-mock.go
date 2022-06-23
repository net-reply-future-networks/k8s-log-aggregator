package mocks

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/sidecar"
)

type MockStream struct {
	MockOpenStderr MockOpenStderr
	MockOpenStdout MockOpenStdout
	MockOpenFile   MockOpenFile
	T              *testing.T
}

type MockOpenStderr struct {
	InPid               sidecar.Pid
	ExpectedInvocations int
	ActualInvocations   int
}

type MockOpenStdout struct {
	InPid               sidecar.Pid
	ExpectedInvocations int
	ActualInvocations   int
}

func (m *MockStream) OpenStderr(pid sidecar.Pid, wg *sync.WaitGroup) {
	defer wg.Done()
	m.MockOpenStderr.ActualInvocations++
	if !reflect.DeepEqual(m.MockOpenStderr.InPid, pid) {
		m.T.Error(fmt.Sprintf("(OpenStderr) unexpected input pid value, expected %v, got %v", m.MockOpenStderr.InPid, pid))
	}
}
func (m *MockStream) OpenStdout(pid sidecar.Pid, wg *sync.WaitGroup) {
	m.MockOpenStdout.ActualInvocations++
	if !reflect.DeepEqual(m.MockOpenStdout.InPid, pid) {
		m.T.Error(fmt.Sprintf("(OpenStdout) unexpected input pid value, expected %v, got %v", m.MockOpenStdout.InPid, pid))
	}
	defer wg.Done()
}
