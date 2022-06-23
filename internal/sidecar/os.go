package sidecar

import "os/exec"

type OsInterface interface {
	ExecPs() ([]byte, error)
}

type Os struct{}

func (o *Os) ExecPs() ([]byte, error) {
	cmd := exec.Command("sh", "-c", `ps -eo pid,comm,ppid | sed 1,1d | awk '{print $1 "," $2 "," $3}' | grep -E -v ',sed,|,ps,|,awk,|,tr,|,sh,|,grep,'`)
	return cmd.Output()
}
