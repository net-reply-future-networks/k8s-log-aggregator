package main

import (
	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/aggregator"
)

func main() {
	ag := new(aggregator.Aggregator)

	ag.Run()
}
