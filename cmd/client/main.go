package main

import (
	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/client"
)

func main() {
	client := client.Client{}
	client.Stream()
}
