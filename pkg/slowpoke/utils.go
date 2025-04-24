package slowpoke

import (
	"golang.org/x/sys/unix"
	"encoding/json"
	"encoding/binary"
	"fmt"
	"net/http"
	"bytes"
	"os"
	"net"
	"time"
	"context"
	"github.com/eniac/mucache/pkg/common"
	"github.com/eniac/mucache/pkg/utility"
	"sync"
	"runtime"
	"sync/atomic"
	// "syscall"
	"google.golang.org/grpc"
	pb "github.com/eniac/mucache/pkg/pb"
	"strings"
)

const printIntervalMillis = 30*1000
// {"thread_id" :{"requestName": counter}}
var (
	// requestCounters map[int]map[string]int
	delayNanos int64
	prerun bool
	requestCounters sync.Map
	accumulatedDelay int64 = 0
	reqCount int64 = 0
	sync_guard sync.RWMutex
	processingMicros int
	sleepSurplus int64 = 0
	pokerBatchThreshold int64

	pipebuf = make([]byte, 8)
	pipefile *os.File
	recover_pipefile *os.File
	grpcConns map[string]*grpc.ClientConn = make(map[string]*grpc.ClientConn)
	grpcConnLock sync.RWMutex
	pipe_recv_buf = make([]byte, 8)
	// grpcConn *grpc.ClientConn
	// initOnce sync.Once

	// for pokerpp
	servName string
	neighbors map[string]net.Conn = make(map[string]net.Conn)
	neighborsLock sync.RWMutex
	isTarget bool
	sleepPhase int
	sleepPhaseLock sync.Mutex
	reqcount = 0
)

func min(a, b int64) int64 {
	if a < b {
			return a
	}
	return b
}

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

func SlowpokeInit() {
	delayMicros := -1
	if env, ok := os.LookupEnv("SLOWPOKE_DELAY_MICROS"); ok {
		fmt.Sscanf(env, "%d", &delayMicros)
		fmt.Printf("SLOWPOKE_DELAY_MICROS=%d\n", delayMicros)
		delayNanos = int64(delayMicros) * int64(1000);
	}
	processingMicros = -1
	if env, ok := os.LookupEnv("SLOWPOKE_PROCESSING_MICROS"); ok {
		fmt.Sscanf(env, "%d", &processingMicros)
		fmt.Printf("SLOWPOKE_PROCESSING_MICROS=%d\n", processingMicros)
	}
	prerun = false
	if env, ok := os.LookupEnv("SLOWPOKE_PRERUN"); ok {
		if env == "true" {
			prerun = true
		}
		fmt.Printf("SLOWPOKE_PRERUN=%t\n", prerun)
	}
	pokerBatchThreshold = 20000000
	if env, ok := os.LookupEnv("SLOWPOKE_POKER_BATCH_THRESHOLD"); ok {
		fmt.Sscanf(env, "%d", &pokerBatchThreshold)
		fmt.Printf("SLOWPOKE_POKER_BATCH_THRESHOLD=%d\n", pokerBatchThreshold)
	}
	if env, ok := os.LookupEnv("SLOWPOKE_SERV_NAME"); ok {
		servName = strings.TrimSpace(env)
		fmt.Printf("SLOWPOKE_SERV_NAME=%s\n", servName)
	}
	isTarget = false
	if env, ok := os.LookupEnv("SLOWPOKE_IS_TARGET_SERVICE"); ok {
		if env == "true" {
			isTarget = true
		} 
		fmt.Printf("env=%s\n", env)
	}
	fmt.Printf("SLOWPOKE_IS_TARGET_SERVICE=%b\n", isTarget)

	var fifo_path string;
	var fifo_recover_path string;
	var ok bool;
	if fifo_path, ok = os.LookupEnv("SLOWPOKE_FIFO_PATH"); ok {
		fmt.Printf("SLOWPOKE_FIFO=%s\n", fifo_path)
	}

	if fifo_recover_path, ok = os.LookupEnv("SLOWPOKE_FIFO_RECOVER_PATH"); ok {
		fmt.Printf("SLOWPOKE_FIFO_RECOVER=%s\n", fifo_recover_path)
	}

	var err error
	pipefile, err = os.OpenFile(fifo_path, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("Error opening pipefile:", err)
		return
	}
	binary.LittleEndian.PutUint64(pipebuf, uint64(0))
	_, err = pipefile.Write(pipebuf);
	if err != nil {
		fmt.Println("Error writing to pipe:", err)
		os.Stdout.Sync()
		return
	}

	recover_pipefile, err = os.OpenFile(fifo_recover_path, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		fmt.Println("Error opening recover pipefile with name :", fifo_recover_path, err)
		return
	}

	_, err = recover_pipefile.Read(pipe_recv_buf);
	for {
		if err == nil {
			break
		}
		fmt.Println("Error reading from pipe:", err)
		os.Stdout.Sync()
		_, err = recover_pipefile.Read(pipe_recv_buf);
	}

	go startControlServer()

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

func SlowpokePokerPPDelay() {
	// Delay
	sync_guard.Lock()
	accumulatedDelay += delayNanos
	reqcount += 1
	if isTarget && int64(reqcount) > pokerBatchThreshold {
		// fmt.Printf("[%s][SlowpokePokerPPDelay] reqcout %d, sleepPhase %d \n", time.Now().Format(time.RFC3339), reqcount, sleepPhase + 1)
		// os.Stdout.Sync()  
		pauseReq := PauseReq{Phase: sleepPhase + 1}
		sync_guard.Unlock()
		handlePause(pauseReq)
		return
	}
	sync_guard.Unlock()
}

func SlowpokeDoDelay(delayToDo int64) {
	// start := time.Now()
	sync_guard.Lock()
	binary.LittleEndian.PutUint64(pipebuf, uint64(delayToDo))
	_, err := pipefile.Write(pipebuf);
	if err != nil {
		fmt.Println("Error writing to pipe:", err)
		os.Stdout.Sync()
		return
	}
	// fmt.Printf("sleeping for: %d [%s]\n", delayToDo, time.Now().Format(time.RFC3339))
	// os.Stdout.Sync()  

	// time.Sleep(time.Duration(accumulatedDelay) * time.Nanosecond)
	_, err = recover_pipefile.Read(pipe_recv_buf);
	if err != nil {
		fmt.Println("Error reading from pipe:", err)
		os.Stdout.Sync()
		panic("pipe")
	}
	sync_guard.Unlock()
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
				break // Exit loop when succeßssful
			}
		
			// Retry in case of failure (another goroutine modified the value)
			counter, _ = requestCounters.Load(serviceFuncName)
		}
	}

	// Process
	lockThread := true
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
			current = getThreadCPUTime();
		}

		takenSurplus += current - target;
		atomic.AddInt64(&sleepSurplus, takenSurplus);

		if lockThread {
			runtime.UnlockOSThread()
		}
	}
}

func SlowpokeDelay() {
	sync_guard.Lock()
	accumulatedDelay += delayNanos
	reqCount++
	// if reqCount >= 100 {
	if accumulatedDelay > pokerBatchThreshold {
		// fmt.Printf("[%s] Should sleep %d \n", time.Now().Format(time.RFC3339), accumulatedDelay)

		binary.LittleEndian.PutUint64(pipebuf, uint64(accumulatedDelay))
		_, err := pipefile.Write(pipebuf);
		if err != nil {
			fmt.Println("Error writing to pipe:", err)
			os.Stdout.Sync()
			return
		}

		// syscall.Syscall(syscall.SYS_SCHED_YIELD, 0, 0, 0) // Yield the CPU

		// elapsed := time.Since(start)
		// start = start.Add(elapsed)
		// accumulatedDelay -= elapsed.Nanoseconds()
		// for accumulatedDelay > 0 {
		//       time.Sleep(time.Duration(accumulatedDelay) * time.Nanosecond)
		//       elapsed = time.Since(start)
		//       start = start.Add(elapsed)
		//       accumulatedDelay -= elapsed.Nanoseconds()
		// }

		// time.Sleep(time.Duration(accumulatedDelay) * time.Nanosecond)
		_, err = recover_pipefile.Read(pipe_recv_buf);
		if err != nil {
			fmt.Println("Error reading from pipe:", err)
			os.Stdout.Sync()
			return
		}
		accumulatedDelay = 0
		reqCount = 0
	}
	sync_guard.Unlock()
}

func SlowpokeFlushDelay(flusher http.Flusher) {
	sync_guard.Lock()
	accumulatedDelay += delayNanos
	reqCount++
	// if accumulatedDelay > pokerBatchThreshold {
	if reqCount >= 100 {
		flusher.Flush()
		
		// start := time.Now()
		binary.LittleEndian.PutUint64(pipebuf, uint64(accumulatedDelay))
		_, err := pipefile.Write(pipebuf);
		if err != nil {
			fmt.Println("Error writing to pipe:", err)
			os.Stdout.Sync()
			return
		}

		// syscall.Syscall(syscall.SYS_SCHED_YIELD, 0, 0, 0) // Yield the CPU

		// elapsed := time.Since(start)
		// start = start.Add(elapsed)
		// accumulatedDelay -= elapsed.Nanoseconds()
		// for accumulatedDelay > 0 {
		//       time.Sleep(time.Duration(accumulatedDelay) * time.Nanosecond)
		//       elapsed = time.Since(start)
		//       start = start.Add(elapsed)
		//       accumulatedDelay -= elapsed.Nanoseconds()
		// }

		// time.Sleep(time.Duration(accumulatedDelay) * time.Nanosecond)
		_, err = recover_pipefile.Read(pipe_recv_buf);
		if err != nil {
			fmt.Println("Error reading from pipe:", err)
			os.Stdout.Sync()
			return
		}
		accumulatedDelay = 0
		reqCount = 0
	}
	sync_guard.Unlock()
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

func Invoke[T interface{}](ctx context.Context, app string, method string, input interface{}) T {
	// sync_guard.RLock()
	// sync_guard.RUnlock()
	requestEndpointRegister(app)
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
	// sync_guard.RLock()
	// sync_guard.RUnlock()
	performRequest[T](ctx, req, &res, app, method, buf)
	// sync_guard.RLock()
	// sync_guard.RUnlock()
	return res
}

func InitGRPCConn(app string) *grpc.ClientConn {
	conn, err := grpc.Dial(fmt.Sprintf("%s.default.svc.cluster.local:%s", app, "80"), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}

func InvokeGRPC(ctx context.Context, app string, method string, input interface{}) string {
	// sync_guard.RLock()
	// sync_guard.RUnlock()

	requestEndpointRegister(app)

	var grpcConn *grpc.ClientConn

	grpcConnLock.RLock()
	if _, ok := grpcConns[app]; ok {
		grpcConn = grpcConns[app]
	}
	grpcConnLock.RUnlock()

	if grpcConn == nil {
		grpcConnLock.Lock()
		if _, ok := grpcConns[app]; !ok {
			grpcConns[app] = InitGRPCConn(app)
		}
		grpcConn = grpcConns[app]
		grpcConnLock.Unlock()
	}

	client := pb.NewSimpleClient(grpcConn)
	resp_, err := client.SimpleRPC(ctx, &pb.SimpleRequest{Endpoint: method})
	if err != nil {
		panic(err)
	}
	return resp_.Resp
}

type SlowpokeListener struct {
	net.Listener
}

func (l *SlowpokeListener) Accept() (net.Conn, error) {
    // sync_guard.RLock()
    // sync_guard.RUnlock()
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
    // sync_guard.RLock()
    // sync_guard.RUnlock()
	return &tracedConn{Conn: conn}, nil
}

type tracedConn struct {
	net.Conn
}

func (c *tracedConn) Read(b []byte) (n int, err error) {
    // sync_guard.RLock()
    // sync_guard.RUnlock()
	n, err = c.Conn.Read(b)
    // sync_guard.RLock()
    // sync_guard.RUnlock()
	return
}

func (c *tracedConn) Write(b []byte) (n int, err error) {
    // sync_guard.RLock()
    // sync_guard.RUnlock()
	n, err = c.Conn.Write(b)
    // sync_guard.RLock()
    // sync_guard.RUnlock()
	return
}

