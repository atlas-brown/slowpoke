package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"golang.org/x/sys/unix"
	"github.com/eniac/mucache/pkg/synthetic"
	// "github.com/eniac/mucache/pkg/slowpoke"
)

var (
	sleepSurplus int64 = 0
)

func getThreadCPUTime() int64 {
	time := unix.Timespec{}
	unix.ClockGettime(unix.CLOCK_THREAD_CPUTIME_ID, &time)
	return time.Nano()
}

func min(a, b int64) int64 {
	if a < b {
			return a
	}
	return b
}

func stressCPU(execTime float32) {
	lockThread := true
	processingMicros := int64(execTime * 1000000.0)
	if processingMicros >= 0 {
		// Threads need to be locked because otherwise util.ThreadCPUTime() can change in the middle of execution
		takenSurplus := atomic.SwapInt64(&sleepSurplus, 0);
		sleepTime := int64(processingMicros*1000.0);
		common := min(takenSurplus, sleepTime);
		takenSurplus -= common;
		sleepTime -= common;
		if lockThread {
			runtime.LockOSThread()
		}

		current := getThreadCPUTime()
		target := current + sleepTime

		for current < target {
			for i := int64(0) ; i < 200000; i++ {
			}
			// slowpoke.RequestRlock()
			current = getThreadCPUTime();
		}

		takenSurplus += current - target;
		atomic.AddInt64(&sleepSurplus, takenSurplus);

		if lockThread {
			runtime.UnlockOSThread()
		}
	}
}

func execCPU(endpoint *synthetic.Endpoint) string {
	complexity := endpoint.CpuComplexity
	threads := complexity.Threads
	executionTime := complexity.ExecutionTime // in seconds

	wg := sync.WaitGroup{}
	wg.Add(threads)

	for i := 0; i < threads; i++ {
		

		go func() {
			defer wg.Done()
			stressCPU(executionTime)
		}()
	}

	wg.Wait()
	return fmt.Sprintf("CPU stress test for %s (%d threads, %f seconds) succeed", endpoint.Name, threads, executionTime)
}

