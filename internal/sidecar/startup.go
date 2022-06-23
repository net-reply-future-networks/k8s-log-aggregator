package sidecar

import "github.com/net-reply-future-networks/k8s-log-aggregator/internal/utils"

func NewSidecar() *Sidecar {
	nc := utils.NewNatsClient()
	lg := NewLogger(nc)
	pm := NewPidManager(lg)
	sm := NewStreamManager(lg)

	sc := Sidecar{
		PidManager:    pm,
		StreamManager: sm,
		Logger:        lg,
	}
	return &sc
}

func NewOsInterface() *Os {
	os := Os{}
	return &os
}

func NewPidManager(lg *Logger) *PidManager {
	pm := PidManager{
		Logger: lg,
		Os:     NewOsInterface(),
	}
	return &pm
}

func NewStreamManager(lg *Logger) *StreamManager {
	sm := StreamManager{
		Logger: lg,
		Os:     NewOsInterface(),
	}
	return &sm
}

func NewLogger(nc *utils.NatsClient) *Logger {
	lg := Logger{
		NatsClient: nc,
	}
	return &lg
}
