package mocks

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/sidecar"
)

type PidManagerMock struct {
	MockExecPs          MockExecPs
	MockGetPids         MockGetPids
	MockConsolidatePids MockConsolidatePids
	T                   *testing.T
}

type MockExecPs struct {
	OutBytes            []byte
	OutErr              error
	ExpectedInvocations int
	ActualInvocations   int
}

type MockGetPids struct {
	OutPids             sidecar.Pids
	OutError            error
	ExpectedInvocations int
	ActualInvocations   int
}

type MockConsolidatePids struct {
	InPids              sidecar.Pids
	OutPids1            sidecar.Pids
	OutPids2            sidecar.Pids
	ExpectedInvocations int
	ActualInvocations   int
}

func (p *PidManagerMock) ExecPs() ([]byte, error) {
	p.MockExecPs.ActualInvocations++
	return p.MockExecPs.OutBytes, p.MockExecPs.OutErr
}
func (p *PidManagerMock) GetPids() (sidecar.Pids, error) {
	p.MockGetPids.ActualInvocations++
	return p.MockGetPids.OutPids, p.MockGetPids.OutError
}
func (p *PidManagerMock) ConsolidatePids(pids sidecar.Pids) (sidecar.Pids, sidecar.Pids) {
	p.MockConsolidatePids.ActualInvocations++
	if !reflect.DeepEqual(pids, p.MockConsolidatePids.InPids) {
		p.T.Error(fmt.Sprintf("(ConsolidatePids) unexpected input pids value, expected %v, got %v", p.MockConsolidatePids.InPids, pids))
	}
	return p.MockConsolidatePids.OutPids1, p.MockConsolidatePids.OutPids2
}
