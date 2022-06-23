package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/sidecar"
)

type TestConsilidatePids struct {
	Test            string
	Pids            sidecar.Pids
	Stored          sidecar.Pids
	ExpectedInPids  sidecar.Pids
	ExpectedOutPids sidecar.Pids
}

func Test_ConsolidatePids(t *testing.T) {
	testTable := []TestConsilidatePids{
		{
			Test: "empty stored",
			Pids: sidecar.Pids{
				{
					Pid:  "1",
					Ppid: "0",
					Comm: "in",
				},
			},
			ExpectedInPids: sidecar.Pids{
				{
					Pid:  "1",
					Ppid: "0",
					Comm: "in",
				},
			},
			ExpectedOutPids: sidecar.Pids{},
			Stored:          sidecar.Pids{},
		},
		{
			Test: "stored matches",
			Pids: sidecar.Pids{
				{
					Pid:  "1",
					Ppid: "0",
					Comm: "in",
				},
			},
			ExpectedInPids:  sidecar.Pids{},
			ExpectedOutPids: sidecar.Pids{},
			Stored: sidecar.Pids{
				{
					Pid:  "1",
					Ppid: "0",
					Comm: "in",
				},
			},
		},
		{
			Test:            "remove stored",
			Pids:            sidecar.Pids{},
			ExpectedInPids:  sidecar.Pids{},
			ExpectedOutPids: sidecar.Pids{},
			Stored: sidecar.Pids{
				{
					Pid:  "1",
					Ppid: "0",
					Comm: "in",
				},
			},
		},

		{
			Test: "replace stored",
			Pids: sidecar.Pids{
				{
					Pid:  "2",
					Ppid: "1",
					Comm: "in",
				},
			},
			ExpectedInPids: sidecar.Pids{
				{
					Pid:  "2",
					Ppid: "1",
					Comm: "in",
				},
			},
			ExpectedOutPids: sidecar.Pids{
				{
					Pid:  "1",
					Ppid: "0",
					Comm: "in",
				},
			},
			Stored: sidecar.Pids{
				{
					Pid:  "1",
					Ppid: "0",
					Comm: "in",
				},
			},
		},
	}

	for _, test := range testTable {
		fmt.Printf("running test %s\n", test.Test)
		pm := sidecar.PidManager{
			Pids: test.Stored,
		}
		out, in := pm.ConsolidatePids(test.Pids)
		if !reflect.DeepEqual(in, test.ExpectedInPids) {
			t.Errorf("unexpected InPids received, expected %v got %v", test.ExpectedInPids, in)
		}
		if !reflect.DeepEqual(in, test.ExpectedInPids) {
			t.Errorf("unexpected OutPids received, expected %v got %v", test.ExpectedOutPids, out)
		}
	}
}
