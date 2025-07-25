package main

import (
	"context"
	"fmt"
	// "github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/common"
	"github.com/eniac/mucache/pkg/wrappers"
	"net/http"
	"runtime"
	"github.com/eniac/mucache/pkg/slowpoke"
	"github.com/eniac/mucache/internal/social"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func readUserTimeline(ctx context.Context, req *social.ReadUserTimelineRequest) *social.ReadUserTimelineResponse {
	// slowpoke.SlowpokeCheck("readUserTimeline")
	posts := social.ReadUserTimeline(ctx, req.UserId)
	//fmt.Printf("Posts read: %+v\n", posts)
	resp := social.ReadUserTimelineResponse{Posts: posts}
	return &resp
}

func writeUserTimeline(ctx context.Context, req *social.WriteUserTimelineRequest) *string {
	// slowpoke.SlowpokeCheck("writeUserTimeline")
	social.WriteUserTimeline(ctx, req.UserId, req.PostIds)
	resp := "OK"
	return &resp
}

func main() {
	if common.ShardEnabled {
		fmt.Println(runtime.GOMAXPROCS(1))
	} else {
		fmt.Println(runtime.GOMAXPROCS(8))
	}
	// go cm.ZmqProxy()
	fmt.Println("Max procs: ", runtime.GOMAXPROCS(0))
	http.HandleFunc("/heartbeat", heartbeat)
	// http.HandleFunc("/ro_read_user_timeline", wrappers.ROWrapper[social.ReadUserTimelineRequest, social.ReadUserTimelineResponse](readUserTimeline))
	http.HandleFunc("/ro_read_user_timeline", wrappers.SlowpokeWrapper[social.ReadUserTimelineRequest, social.ReadUserTimelineResponse](readUserTimeline, "ro_read_user_timeline"))
	// http.HandleFunc("/write_user_timeline", wrappers.NonROWrapper[social.WriteUserTimelineRequest, string](writeUserTimeline))
	http.HandleFunc("/write_user_timeline", wrappers.SlowpokeWrapper[social.WriteUserTimelineRequest, string](writeUserTimeline, "write_user_timeline"))
	slowpoke.SlowpokeInit()
	fmt.Println("Starting server on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
