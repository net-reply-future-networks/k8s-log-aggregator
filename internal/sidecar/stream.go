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
	Logger  LoggerInterface
	Os      OsInterface
}

type StreamManagerInterface interface {
	OpenStream(pid Pid)
	CloseStream(pid Pid) bool
}

type Stream struct {
	Pid          Pid
	CancelStdout chan bool
	CancelStderr chan bool
	Logger       LoggerInterface
	Os           OsInterface
}

type StreamInterface interface {
	OpenStderr(pid Pid, wg *sync.WaitGroup)
	OpenStdout(pid Pid, wg *sync.WaitGroup)
	OpenFile(name string) (*os.File, error)
}

func (s *Stream) OpenStdout(pid Pid, wg *sync.WaitGroup) {
	defer wg.Done()
	s.Logger.InfoPID("Stdout stream opening", pid)
	f, err := s.Os.OpenFile(fmt.Sprintf("/proc/%s/fd/1", pid.Pid))
	if err != nil {
		s.Logger.ErrorPID(err.Error(), pid)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		select {
		case <-s.CancelStdout:
			s.Logger.InfoPID("Stdout stream closed", pid)
			return
		default:
			line, err := r.ReadBytes('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				s.Logger.ErrorPID(err.Error(), pid)
				continue
			} else if err != nil && errors.Is(err, io.EOF) {
				continue
			}
			if len(line) == 0 {
				continue
			}
			s.Logger.InfoPID(string(line), pid)
		}
	}
}

func (s *Stream) OpenStderr(pid Pid, wg *sync.WaitGroup) {
	defer wg.Done()
	s.Logger.InfoPID("Stderr stream opening", pid)
	f, err := s.Os.OpenFile(fmt.Sprintf("/proc/%s/fd/2", pid.Pid))
	if err != nil {
		s.Logger.ErrorPID(err.Error(), pid)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		select {
		case <-s.CancelStderr:
			s.Logger.InfoPID("Stderr stream closed", pid)
			return
		default:
			line, err := r.ReadBytes('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				s.Logger.ErrorPID(err.Error(), pid)
				continue
			} else if err != nil && errors.Is(err, io.EOF) {
				continue
			}
			if len(line) == 0 {
				continue
			}
			s.Logger.ErrorPID(string(line), pid)
		}
	}
}

func (sm *StreamManager) OpenStream(pid Pid) {
	stream := new(Stream)
	stream.Pid = pid
	stream.Os = sm.Os
	stream.Logger = sm.Logger
	stream.CancelStdout = make(chan bool, 1)
	stream.CancelStderr = make(chan bool, 1)
	sm.wg.Add(2)
	go stream.OpenStdout(pid, &sm.wg)
	go stream.OpenStderr(pid, &sm.wg)
	sm.Streams = append(sm.Streams, stream)
}

func (sm *StreamManager) CloseStream(pid Pid) bool {
	sm.Logger.InfoPID("Stream Closing", pid)
	for i, stream := range sm.Streams {
		if stream.Pid.Pid == pid.Pid {
			stream.CancelStdout <- true
			stream.CancelStderr <- true
			sm.Streams = append(sm.Streams[:i], sm.Streams[i+1:]...)
			return true
		}
	}
	return false
}
