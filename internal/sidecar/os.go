package sidecar

import (
	"os"
	"os/exec"
)

type OsInterface interface {
	ExecPs() ([]byte, error)
	OpenFile(name string) (*os.File, error)
}

type Os struct{}

func (o *Os) ExecPs() ([]byte, error) {
	cmd := exec.Command("sh", "-c", `ps -eo pid,comm,ppid | sed 1,1d | awk '{print $1 "," $2 "," $3}' | grep -E -v ',sed,|,ps,|,awk,|,tr,|,sh,|,grep,'`)
	return cmd.Output()
}

func (o *Os) OpenFile(name string) (*os.File, error) {
	return os.Open(name)
}
