package main

import (
	"golang.org/x/sys/unix"
	"os"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"net/http"
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

func MoreComplicatedSlowpoke() {
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
	takenSurplus := atomic.SwapInt64(&sleepSurplus, 0);
	sleepTime := int64(delayMicros*1000.0);
	common := min(takenSurplus, sleepTime);
	takenSurplus -= common;
	sleepTime -= common;
	runtime.LockOSThread()
	current := getThreadCPUTime()
	target := current + sleepTime

	for current < target {
		current = getThreadCPUTime();
	}

	takenSurplus += current - target;
	atomic.AddInt64(&sleepSurplus, takenSurplus);

	runtime.UnlockOSThread()
}

func SimpleSpin() {
     for i := 0; i < delayMicros; i++ {
	     for j := 0; j < 1000; j++ {

	     }
     }
}

func SetupWorkers(n int, kind string, wg *sync.WaitGroup) {
	if kind == "spinx" {
		for w := 0; w < n; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 200; i++ {
					ComplicatedSlowpoke()
				}
			}()
		}
	} else if kind == "simple" {
		for w := 0; w < n; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for i := 0; i < 200; i++ {
					SimpleSpin();
				}
			}()
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	ComplicatedSlowpoke();
	for i := 0; i < 80000; i++ {}
	fmt.Fprintf(w, "Welcome to the Home Page!")
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
	http.HandleFunc("/", homeHandler)
	fmt.Println("Server is running on http://localhost:8000", kind, delayMicros)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
	//now := time.Now()
        //nano := now.UnixNano()
	//SetupWorkers(20, kind, &wg)
	//wg.Wait()
	//now = time.Now()
	//done_nano := now.UnixNano()
        //fmt.Printf("%d\n", (done_nano - nano)/1000);
}
