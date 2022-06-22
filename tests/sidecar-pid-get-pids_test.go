package tests

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/sidecar"
	"github.com/net-reply-future-networks/k8s-log-aggregator/tests/mocks"
)

type TestGetPidsInput struct {
	Test          string
	MockExecPs    mocks.MockExecPs
	ExpectedPids  sidecar.Pids
	ExpectedError error
}

func Test_GetPids(t *testing.T) {
	testTable := []TestGetPidsInput{
		{
			Test: "happy path",
			MockExecPs: mocks.MockExecPs{
				ExpectedInvocations: 1,
				OutBytes:            []byte("\n1,wrong\n1,correct,0\nwrong\n/\n"),
				OutErr:              nil,
			},
			ExpectedPids: sidecar.Pids{
				{
					Pid:  "1",
					Comm: "correct",
					Ppid: "0",
				},
			},
			ExpectedError: nil,
		},
		{
			Test: "ExecPs fails",
			MockExecPs: mocks.MockExecPs{
				ExpectedInvocations: 1,
				OutBytes:            nil,
				OutErr:              errors.New("test error"),
			},
			ExpectedPids:  sidecar.Pids{},
			ExpectedError: errors.New("test error"),
		},
	}

	for _, test := range testTable {
		fmt.Printf("running test %s\n", test.Test)
		mock := mocks.OsMock{
			MockExecPs: &test.MockExecPs,
		}
		pm := sidecar.PidManager{
			Os: &mock,
		}

		pids, err := pm.GetPids()

		if !reflect.DeepEqual(pids, test.ExpectedPids) {
			t.Errorf("unexpected pids received, expected %v got %v", test.ExpectedPids, pids)
		}
		if (err != nil && test.ExpectedError != nil) && errors.Is(test.ExpectedError, err) {
			t.Errorf("unexpected error received, expected %v, got %v", test.ExpectedError, err)
		}
		if test.MockExecPs.ExpectedInvocations != test.MockExecPs.ActualInvocations {
			t.Errorf("unexpected invocations of MockExecPs, expected %d, got %d", test.MockExecPs.ExpectedInvocations, test.MockExecPs.ActualInvocations)
		}
	}
}
