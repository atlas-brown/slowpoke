package slowpoke

import (
	"golang.org/x/sys/unix"
	"encoding/json"
	"encoding/binary"
	"fmt"
	"net/http"
	"bytes"
	"os"
	"time"
	"context"
	"github.com/eniac/mucache/pkg/common"
	"github.com/eniac/mucache/pkg/utility"
	"sync"
	"sync/atomic"
)

const printIntervalMillis = 30*1000
// {"thread_id" :{"requestName": counter}}
var (
	// requestCounters map[int]map[string]int
	delayMicros int
	prerun bool
	requestCounters sync.Map
	sleepSurplus int64 = 0

	req_events int64
)

func min(a, b int64) int64 {
	if a < b {
			return a
	}
	return b
}
// func printCounters() {
// 	time_now := time.Now()
// 	func_counters := make(map[string]int)
// 	string_to_print := ""
// 	for thread, counters := range requestCounters {
// 		for requestName, count := range counters {
// 			if count > 0 {
// 				if _, ok := func_counters[requestName]; !ok {
// 					func_counters[requestName] = 0
// 				}
// 				func_counters[requestName] += count
// 				string_to_print += fmt.Sprintf("	[%d] %s: %d\n", thread, requestName, count)
// 				requestCounters[thread][requestName] = 0
// 			}
// 		}
// 	}
// 	if string_to_print != "" {
// 		fmt.Printf("[%s] Slowpoke Counters\n", time_now.String())
// 		fmt.Printf(string_to_print)
// 		fmt.Printf("	[Aggregation]\n")
// 		for func_name, count := range func_counters {
// 			fmt.Printf("		%s: %d\n", func_name, count)
// 		}
// 	}
// }

func printCountersSyncMap() {
	timeNow := time.Now()
	funcCounters := make(map[string]int)
	totalCounter := 0

	// Iterate over sync.Map safely
	requestCounters.Range(func(funcName, counter interface{}) bool {
		funcNameStr, ok1 := funcName.(string)
		counterInt, ok2 := counter.(int)

		if !ok1 || !ok2 {
			fmt.Println("Invalid type in sync.Map") // Prevent panic
			return true
		}

		// Ensure no overwrites: Reset counter **before** adding it to funcCounters
		for !requestCounters.CompareAndSwap(funcNameStr, counterInt, 0) {
			// Reload latest value to prevent overwriting concurrent increments
			newCounter, _ := requestCounters.Load(funcNameStr)
			counterInt, _ = newCounter.(int)
		}

		// Now increment local counters after resetting the global counter
		if counterInt > 0 {
			funcCounters[funcNameStr] += counterInt
			totalCounter += counterInt
		}

		return true
	})

	// Print results
	msg := ""
	msg += fmt.Sprintf("[%s] Slowpoke Counters\n", timeNow.Format(time.RFC3339))
	for funcName, count := range funcCounters {
		msg += fmt.Sprintf("\t%s: %d\n", funcName, count)
	}
	msg += fmt.Sprintf("\t[Total] %d\n", totalCounter)
	if totalCounter > 0 {
		fmt.Println(msg)
	}
}

// Saves the response to *res (also might save the result to cache if we are in upperbound baseline
func performRequest[T interface{}](ctx context.Context, req *http.Request, res *T, app string, method string, argBytes []byte) {
	resp, err := common.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	utility.Assert(resp.StatusCode == http.StatusOK)
	defer resp.Body.Close()
	utility.ParseJson(resp.Body, res)
}

func SlowpokeInit() {
	delayMicros = -1
	if env, ok := os.LookupEnv("SLOWPOKE_DELAY_MICROS"); ok {
		fmt.Sscanf(env, "%d", &delayMicros)
		fmt.Printf("SLOWPOKE_DELAY_MICROS=%d\n", delayMicros)
	}
	prerun = false
	if env, ok := os.LookupEnv("SLOWPOKE_PRERUN"); ok {
		if env == "true" {
			prerun = true
		}
		fmt.Printf("SLOWPOKE_PRERUN=%t\n", prerun)
	}

	var fifo_path string;
	var ok bool;
	if fifo_path, ok = os.LookupEnv("SLOWPOKE_FIFO_PATH"); ok {
		fmt.Printf("SLOWPOKE_FIFO=%s\n", fifo_path)
	}

	fmt.Println("wtf:")
	go func () {
		file, err := os.OpenFile(fifo_path, os.O_WRONLY, os.ModeNamedPipe)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()
		prev_i := int64(0)
		i := req_events
		buf := make([]byte, 8)
		for {
			<-time.After(10 * time.Millisecond)
			i = req_events
			if i - prev_i > 5000 {
				value := (i - prev_i) * int64(delayMicros) * int64(1000)
				binary.LittleEndian.PutUint64(buf, uint64(value))
				_, err = file.Write(buf);
				if err != nil {
					fmt.Println("Error writing to pipe:", err)
					os.Stdout.Sync()
					return
				}
			    	/* fmt.Println("get more than 10000 requests"); */
				prev_i = i
			}
		}
	}()

	if !prerun {
		return
	}

	// requestCounters = make(map[int]map[string]int)
	requestCounters = sync.Map{}
	go func() {
		for {
			<-time.After(printIntervalMillis * time.Millisecond)
			// printCounters()
			printCountersSyncMap()
		}
	}()
}

// Get the amount of time in nanoseconds the calling thread has spent using the CPU since startup
func getThreadCPUTime() int64 {
	time := unix.Timespec{}
	unix.ClockGettime(unix.CLOCK_THREAD_CPUTIME_ID, &time)
	return time.Nano()
}


func SlowpokeCheck(serviceFuncName string) {
	// // Record request
	// if _, ok := requestCounters[unix.Gettid()]; !ok {
	// 	requestCounters[unix.Gettid()] = make(map[string]int)
	// }
	// if _, ok := requestCounters[unix.Gettid()][serviceFuncName]; !ok {
	// 	requestCounters[unix.Gettid()][serviceFuncName] = 0
	// }
	// requestCounters[unix.Gettid()][serviceFuncName]++

	if prerun{
		// Record request
		counter, _ := requestCounters.LoadOrStore(serviceFuncName, 0)

		for {
			current := counter.(int) // Safely type assert
			if requestCounters.CompareAndSwap(serviceFuncName, current, current+1) {
				break // Exit loop when successful
			}
		
			// Retry in case of failure (another goroutine modified the value)
			counter, _ = requestCounters.Load(serviceFuncName)
		}
	}

	// Delay
	atomic.AddInt64(&req_events, 1);
}

func Invoke[T interface{}](ctx context.Context, app string, method string, input interface{}) T {
	buf, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	var res T
	// Use kubernete native DNS addr
	url := fmt.Sprintf("http://%s.%s.svc.cluster.local:%s/%s", app, "default", "80", method)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	if err != nil {
		panic(err)
	}
	performRequest[T](ctx, req, &res, app, method, buf)
	return res
}
