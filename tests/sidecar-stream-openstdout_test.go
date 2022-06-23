package tests

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/sidecar"
	"github.com/net-reply-future-networks/k8s-log-aggregator/tests/mocks"
)

type TestOpenStdOut struct {
	Pid              sidecar.Pid
	InfoPidExpected  []string
	ErrorPidExpected []string
	FileWrites       []string
	OpenFileOutErr   error
}

func Test_OpenStdOut(t *testing.T) {
	testTable := []TestOpenStdOut{
		{
			Pid: sidecar.Pid{
				Pid: "1",
			},
			FileWrites: []string{
				"test log\n",
				"test log 2",
			},
			OpenFileOutErr: nil,
			InfoPidExpected: []string{
				"Stdout stream opening",
				"test log\n",
				"test log 2\n",
				"Stdout stream closed",
			},
		},
	}

	for _, test := range testTable {
		r, w, _ := os.Pipe()
		om := mocks.OsMock{
			MockOpenFile: mocks.MockOpenFile{
				OutErr:  test.OpenFileOutErr,
				OutFile: r,
			},
		}
		lg := mocks.LogMock{}
		stream := sidecar.Stream{
			Os:           &om,
			CancelStdout: make(chan bool, 1),
			Logger:       &lg,
			Pid:          test.Pid,
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go stream.OpenStdout(test.Pid, &wg)

		for _, x := range test.FileWrites {
			if _, err := w.Write([]byte(x)); err != nil {
				t.Error(err)
			}
		}
		time.Sleep(1 * time.Second)
		stream.CancelStdout <- true
		if _, err := w.Write([]byte("\n")); err != nil {
			t.Error(err)
		}
		wg.Wait()
		if !reflect.DeepEqual(lg.ErrorPIDOut, test.ErrorPidExpected) {
			t.Errorf("unexpected errorPID log output, expected %v, got %v\n", test.ErrorPidExpected, lg.ErrorPIDOut)
		}
		if !reflect.DeepEqual(lg.InfoPIDOut, test.InfoPidExpected) {
			t.Errorf("unexpected errorPID log output, expected %v, got %v\n", test.InfoPidExpected, lg.InfoPIDOut)
		}
	}

}
