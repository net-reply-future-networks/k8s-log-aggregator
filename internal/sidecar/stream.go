package sidecar

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

type StreamManager struct {
	Streams []*Stream
	wg      sync.WaitGroup
}

type Stream struct {
	Pid    Pid
	Cancel chan bool
}

func (s *Stream) Open(pid Pid, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("--- Stream Opening ---\nPID:  %s\nPPID: %s\nCOMM: %s\n----------------------\n", pid.Pid, pid.Ppid, pid.Comm)
	f, err := os.Open(fmt.Sprintf("/proc/%s/fd/1", pid.Pid))
	if err != nil {
		fmt.Println(err)
		return
	}
	r := bufio.NewReader(f)
	for {
		select {
		case <-s.Cancel:
			fmt.Printf("PID %s stream closed\n", pid.Pid)
			return
		default:
			line, err := r.ReadBytes('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				fmt.Printf("pid: %s - err: %s\n", pid.Pid, err.Error())
				continue
			}
			if len(line) == 0 {
				continue
			}
			fmt.Println(line)
			fmt.Printf("pid: %s - log: %s\n", pid.Pid, string(line))
		}
	}
}

func (sm *StreamManager) OpenStream(pid Pid) {
	stream := new(Stream)
	stream.Pid = pid
	stream.Cancel = make(chan bool, 1)
	sm.wg.Add(1)
	go stream.Open(pid, &sm.wg)
	sm.Streams = append(sm.Streams, stream)
}

func (sm *StreamManager) CloseStream(pid Pid) bool {
	fmt.Println("Closing stream", pid)
	for i, stream := range sm.Streams {
		if stream.Pid.Pid == pid.Pid {
			stream.Cancel <- true
			sm.Streams = append(sm.Streams[:i], sm.Streams[i+1:]...)
			return true
		}
	}
	fmt.Println(sm.Streams)
	return false
}
