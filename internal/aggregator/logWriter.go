package aggregator

import (
	"encoding/json"
	"fmt"

	"github.com/net-reply-future-networks/k8s-log-aggregator/internal/sidecar"
)

type LogWriter struct {
	DbHandlers DbHandlers
}

func (lw *LogWriter) WriteLog(b []byte) {
	log := sidecar.Log{}
	if err := json.Unmarshal(b, &log); err != nil {
		fmt.Println(err)
	}
	token := lw.DbHandlers.GenerateKey()

	fmt.Printf("\nStoring log:\n%s\nKey:%s\n\n", string(b), string(token))

	lw.DbHandlers.SetRecord(string(token), log)
}
