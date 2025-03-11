package main

import (
	"golang.org/x/sys/unix"
	"os"
	"fmt"
	"time"
	"runtime"
	"syscall"
	"strconv"
	"sync/atomic"
)

var (
	delayMicros int
	sleepSurplus int64 = 0
)

func min(a, b int64) int64 {
	if a < b {
			return a
	}
	return b
}

func getThreadCPUTime() int64 {
	time := unix.Timespec{}
	unix.ClockGettime(unix.CLOCK_THREAD_CPUTIME_ID, &time)
	return time.Nano()
}

func ComplicatedSlowpoke() {
	takenSurplus := atomic.SwapInt64(&sleepSurplus, 0);
	sleepTime := int64(delayMicros*1000.0);
	common := min(takenSurplus, sleepTime);
	takenSurplus -= common;
	sleepTime -= common;
	runtime.LockOSThread()
	current := getThreadCPUTime()
	target := current + sleepTime

	for current < target {
		for i := int64(0) ; i < 200000; i++ {
		}
		current = getThreadCPUTime();
	}

	takenSurplus += current - target;
	atomic.AddInt64(&sleepSurplus, takenSurplus);

	runtime.UnlockOSThread()
}

func SlowpokeCheck() {
	lockThread := true
	if delayMicros >= 0 {
		if lockThread {
			runtime.LockOSThread()
		}

		start := getThreadCPUTime()
		target := start + int64(delayMicros*1000.0)

		for getThreadCPUTime() < target {
		}

		if lockThread {
			runtime.UnlockOSThread()
		}
	}
}

func SimpleSpin() {
     for i := 0; i < delayMicros; i++ {
           syscall.Syscall(syscall.SYS_GETPID, 0, 0, 0);
     }
}

func main() {
	args := os.Args
	var kind string
	if len(args) > 2 {
		num, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Invalid number:", err)
			return
		}
		delayMicros = num
		kind = args[1]
	} else {
		fmt.Println("wrong args")
		return
	}
	now := time.Now()
        nano := now.UnixNano()
	if kind == "getpid" {
		for i := 0; i < 2000; i++ {
			SimpleSpin()
        	}
	} else if kind == "spinx" {
		for i := 0; i < 2000; i++ {
			ComplicatedSlowpoke();
        	}
	} else if kind == "spin" {
		for i := 0; i < 2000; i++ {
			SlowpokeCheck();
        	}
	} else {
		fmt.Println("wrong args")
		return
	}
	now = time.Now()
	done_nano := now.UnixNano()
        fmt.Printf("%d\n", (done_nano - nano)/1000);
}
