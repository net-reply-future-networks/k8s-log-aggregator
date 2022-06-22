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
	Logger  Logger
}

type Stream struct {
	Pid    Pid
	Cancel chan bool
	Logger Logger
}

func (s *Stream) OpenFile(name string) (*os.File, error) {
	return os.Open(name)
}

func (s *Stream) Open(pid Pid, wg *sync.WaitGroup) {
	defer wg.Done()
	s.Logger.InfoPID("Stream opening", pid)
	f, err := s.OpenFile(fmt.Sprintf("/proc/%s/fd/1", pid.Pid))
	if err != nil {
		s.Logger.ErrorPID(err.Error(), pid)
		return
	}
	r := bufio.NewReader(f)
	for {
		select {
		case <-s.Cancel:
			s.Logger.InfoPID("Stream Closed", pid)
			return
		default:
			line, err := r.ReadBytes('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				s.Logger.ErrorPID(err.Error(), pid)
				continue
			}
			if len(line) == 0 {
				continue
			}
			s.Logger.InfoPID(string(line), pid)
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
	sm.Logger.InfoPID("Stream Closing", pid)
	for i, stream := range sm.Streams {
		if stream.Pid.Pid == pid.Pid {
			stream.Cancel <- true
			sm.Streams = append(sm.Streams[:i], sm.Streams[i+1:]...)
			return true
		}
	}
	return false
}
