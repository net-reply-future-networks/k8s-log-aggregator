package sidecar

import (
	"fmt"
	"log"
	"time"
)

type Sidecar struct {
	PidManager    PidManager
	StreamManager StreamManager
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
				fmt.Println("Unable to find pid, possible run away go routine", pid)
			}
		}
		for _, pid := range newPids {
			s.StreamManager.OpenStream(pid)
		}

		time.Sleep(2 * time.Second)
	}
}
