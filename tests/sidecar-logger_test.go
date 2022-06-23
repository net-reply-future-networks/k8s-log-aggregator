package tests

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/sidecar"
	"github.com/net-reply-future-networks/k8s-log-aggregator/tests/mocks"
)

type TestLog struct {
	Test               string
	Input              string
	ExpectedLog        sidecar.Log
	ExpectedSidecarLog []byte
	Pid                sidecar.Pid
	MockPublish        mocks.MockPublish
}

func Test_LogError(t *testing.T) {
	testTable := []TestLog{
		{
			Test:  "happy path",
			Input: "im a log",
			ExpectedLog: sidecar.Log{
				Log:   "im a log",
				Level: "error",
			},
			MockPublish: mocks.MockPublish{
				OutError:            nil,
				ExpectedInvocations: 1,
			},
		},
	}
	for _, test := range testTable {
		fmt.Printf("running test %s\n", test.Test)
		nc := mocks.NatsClientMock{
			MockPublish: test.MockPublish,
			T:           t,
		}
		lg := sidecar.Logger{
			NatsClient: &nc,
		}
		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}
		os.Stdout = w
		lg.Error(test.Input)
		os.Stdout = old
		reader := bufio.NewReader(r)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			t.Error(err)
		}
		log := sidecar.Log{}
		err = json.Unmarshal(line, &log)

		if log.Level != test.ExpectedLog.Level {
			t.Errorf("unexpected Log Level received, expected %s got %s", test.ExpectedLog.Level, log.Level)
		}
		if log.Log != test.ExpectedLog.Log {
			t.Errorf("unexpected Log Log received, expected %s got %s", test.ExpectedLog.Log, log.Log)
		}
	}
}

func Test_LogErrorPid(t *testing.T) {
	testTable := []TestLog{
		{
			Test:  "happy path",
			Input: "im a log",
			ExpectedLog: sidecar.Log{
				Log:   "im a log",
				Level: "error",
				Process: &sidecar.Pid{
					Pid:  "1",
					Ppid: "2",
					Comm: "3",
				},
			},
			MockPublish: mocks.MockPublish{
				OutError:            nil,
				ExpectedInvocations: 1,
			},
			Pid: sidecar.Pid{
				Pid:  "1",
				Ppid: "2",
				Comm: "3",
			},
		},
	}
	for _, test := range testTable {
		fmt.Printf("running test %s\n", test.Test)
		nc := mocks.NatsClientMock{
			MockPublish: test.MockPublish,
			T:           t,
		}
		lg := sidecar.Logger{
			NatsClient: &nc,
		}
		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}
		os.Stdout = w
		lg.ErrorPID(test.Input, test.Pid)
		os.Stdout = old
		reader := bufio.NewReader(r)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			t.Error(err)
		}
		log := sidecar.Log{}
		err = json.Unmarshal(line, &log)

		if log.Level != test.ExpectedLog.Level {
			t.Errorf("unexpected Log Level received, expected %s got %s", test.ExpectedLog.Level, log.Level)
		}
		if log.Log != test.ExpectedLog.Log {
			t.Errorf("unexpected Log Log received, expected %s got %s", test.ExpectedLog.Log, log.Log)
		}
		if reflect.DeepEqual(log.Process, test.Pid) {
			t.Errorf("unexpected Log Pid received, expected %v got %v", test.Pid, log.Process)
		}
	}
}

func Test_InfoError(t *testing.T) {
	testTable := []TestLog{
		{
			Test:  "happy path",
			Input: "im a log",
			ExpectedLog: sidecar.Log{
				Log:   "im a log",
				Level: "info",
			},
			MockPublish: mocks.MockPublish{
				OutError:            nil,
				ExpectedInvocations: 1,
			},
		},
	}
	for _, test := range testTable {
		fmt.Printf("running test %s\n", test.Test)
		nc := mocks.NatsClientMock{
			MockPublish: test.MockPublish,
			T:           t,
		}
		lg := sidecar.Logger{
			NatsClient: &nc,
		}
		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}
		os.Stdout = w
		lg.Info(test.Input)
		os.Stdout = old
		reader := bufio.NewReader(r)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			t.Error(err)
		}
		log := sidecar.Log{}
		err = json.Unmarshal(line, &log)

		if log.Level != test.ExpectedLog.Level {
			t.Errorf("unexpected Log Level received, expected %s got %s", test.ExpectedLog.Level, log.Level)
		}
		if log.Log != test.ExpectedLog.Log {
			t.Errorf("unexpected Log Log received, expected %s got %s", test.ExpectedLog.Log, log.Log)
		}
	}
}

func Test_LogInfoPid(t *testing.T) {
	testTable := []TestLog{
		{
			Test:  "happy path",
			Input: "im a log",
			ExpectedLog: sidecar.Log{
				Log:   "im a log",
				Level: "info",
				Process: &sidecar.Pid{
					Pid:  "1",
					Ppid: "2",
					Comm: "3",
				},
			},
			MockPublish: mocks.MockPublish{
				OutError:            nil,
				ExpectedInvocations: 1,
			},
			Pid: sidecar.Pid{
				Pid:  "1",
				Ppid: "2",
				Comm: "3",
			},
		},
	}
	for _, test := range testTable {
		fmt.Printf("running test %s\n", test.Test)
		nc := mocks.NatsClientMock{
			MockPublish: test.MockPublish,
			T:           t,
		}
		lg := sidecar.Logger{
			NatsClient: &nc,
		}
		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}
		os.Stdout = w
		lg.InfoPID(test.Input, test.Pid)
		os.Stdout = old
		reader := bufio.NewReader(r)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			t.Error(err)
		}
		log := sidecar.Log{}
		err = json.Unmarshal(line, &log)

		if log.Level != test.ExpectedLog.Level {
			t.Errorf("unexpected Log Level received, expected %s got %s", test.ExpectedLog.Level, log.Level)
		}
		if log.Log != test.ExpectedLog.Log {
			t.Errorf("unexpected Log Log received, expected %s got %s", test.ExpectedLog.Log, log.Log)
		}
		if reflect.DeepEqual(log.Process, test.Pid) {
			t.Errorf("unexpected Log Pid received, expected %v got %v", test.Pid, log.Process)
		}
	}
}

func Test_SidecarLog(t *testing.T) {
	testTable := []TestLog{
		{
			Test:               "happy path",
			Input:              "im a log",
			ExpectedSidecarLog: []byte("im a log\n"),
		},
	}
	for _, test := range testTable {
		fmt.Printf("running test %s\n", test.Test)
		lg := sidecar.Logger{}
		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			t.Error(err)
		}
		os.Stdout = w
		lg.SidecarLog(test.Input)
		os.Stdout = old
		reader := bufio.NewReader(r)
		line, err := reader.ReadBytes('\n')
		if err != nil {
			t.Error(err)
		}
		if string(line) != string(test.ExpectedSidecarLog) {
			t.Errorf("unexpected Log Log received, expected %s got %s", string(line), test.ExpectedSidecarLog)
		}
	}
}
