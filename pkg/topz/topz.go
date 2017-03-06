package topz

import (
	"net/http"

	"time"

	"log"

	"sync"

	"encoding/json"

	"github.com/shirou/gopsutil/process"
)

func handleError(res http.ResponseWriter, err error) {
	res.WriteHeader(http.StatusInternalServerError)
	res.Write([]byte(err.Error()))
}

type ProcInfo struct {
	PID           int32
	MemoryPercent float32
	CPUPercent    float64
}

func HandleTopz(res http.ResponseWriter, req *http.Request) {
	pids, err := process.Pids()
	if err != nil {
		handleError(res, err)
		return
	}
	procs := make([]ProcInfo, len(pids))
	wg := sync.WaitGroup{}
	wg.Add(len(pids))
	for ix := range pids {
		proc, err := process.NewProcess(pids[ix])
		if err != nil {
			handleError(res, err)
			return
		}
		go func(i int) {
			var err error
			procs[i].PID = pids[i]
			if procs[i].MemoryPercent, err = proc.MemoryPercent(); err != nil {
				log.Printf("Error getting Memory: %v", err)
			}
			if procs[i].CPUPercent, err = proc.Percent(100 * time.Millisecond); err != nil {
				log.Printf("Error getting CPU: %v", err)
			}
			wg.Done()
		}(ix)
	}
	wg.Wait()

	res.WriteHeader(http.StatusOK)
	b, err := json.Marshal(procs)
	if err != nil {
		handleError(res, err)
		return
	}
	res.Write(b)
}
