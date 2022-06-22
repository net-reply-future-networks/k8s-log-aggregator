package sidecar

import (
	"log"
	"time"
)

type Sidecar struct {
	PidManager    PidManagerInterface
	StreamManager StreamManagerInterface
	Logger        LoggerInterface
}

func (s *Sidecar) Run() {
	for {
		pids, err := s.PidManager.GetPids()
		if err != nil {
			log.Fatal(err)
		}
		oldPids, newPids := s.PidManager.ConsolidatePids(pids)

		for _, pid := range oldPids {
			if ok := s.StreamManager.CloseStream(pid); !ok {
				s.Logger.ErrorPID("Unable to find pid, possible run away go routine", pid)
			}
		}
		for _, pid := range newPids {
			s.StreamManager.OpenStream(pid)
		}

		time.Sleep(10 * time.Second)
	}
}
