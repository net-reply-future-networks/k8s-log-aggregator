package sidecar

import (
	"fmt"
	"os"
	"strings"
)

type PidManager struct {
	Pids   Pids
	Logger LoggerInterface
	Os     OsInterface
}

type PidManagerInterface interface {
	GetPids() (Pids, error)
	ConsolidatePids(pids Pids) (Pids, Pids)
}

type Pid struct {
	Pid  string `json:"pid"`
	Ppid string `json:"ppid"`
	Comm string `json:"comm"`
}

type Pids []Pid

func (p Pids) Contains(pid string) bool {
	for _, x := range p {
		if x.Pid == pid {
			return true
		}
	}
	return false
}

func (p *PidManager) GetPids() (Pids, error) {
	myPid := fmt.Sprintf("%d", os.Getpid())
	pids := Pids{}

	out, err := p.Os.ExecPs()
	if err != nil {
		return pids, err
	}

	rows := strings.Split(string(out), "\n")
	for _, x := range rows {
		split := strings.Split(x, ",")
		if len(split) != 3 {
			continue
		}
		if split[0] == myPid {
			continue
		}
		pids = append(pids, Pid{
			Pid:  split[0],
			Ppid: split[2],
			Comm: split[1],
		})
	}
	return pids, nil
}

func (p *PidManager) ConsolidatePids(pids Pids) (Pids, Pids) {
	newPids := Pids{}
	oldPids := Pids{}

	for _, pid := range pids {
		if !p.Pids.Contains(pid.Pid) {
			newPids = append(newPids, pid)
		}
	}
	for _, pid := range p.Pids {
		if !pids.Contains(pid.Pid) {
			oldPids = append(oldPids, pid)
		}
	}

	p.Pids = pids

	return oldPids, newPids
}
