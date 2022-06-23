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

type TestOpenStderr struct {
	Pid              sidecar.Pid
	InfoPidExpected  []string
	ErrorPidExpected []string
	FileWrites       []string
	OpenFileOutErr   error
}

func Test_OpenStdErr(t *testing.T) {
	go func() {
		time.Sleep(5 * time.Second)
		t.Error("TIMED OUT")
	}()
	testTable := []TestOpenStderr{
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
				"Stderr stream opening",
				"Stderr stream closed",
			},
			ErrorPidExpected: []string{
				"test log\n",
				"test log 2\n",
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
			CancelStderr: make(chan bool, 1),
			Logger:       &lg,
			Pid:          test.Pid,
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go stream.OpenStderr(test.Pid, &wg)

		for _, x := range test.FileWrites {
			if _, err := w.Write([]byte(x)); err != nil {
				t.Error(err)
			}
		}
		time.Sleep(1 * time.Second)
		stream.CancelStderr <- true
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
