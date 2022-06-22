package main

import (
	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/sidecar"
)

func main() {
	sc := sidecar.NewSidecar()

	sc.Run()
}
