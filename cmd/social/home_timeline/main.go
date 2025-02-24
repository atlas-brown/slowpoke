package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/social"
	// "github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/common"
	"github.com/eniac/mucache/pkg/wrappers"
	"github.com/eniac/mucache/pkg/slowpoke"
	"net/http"
	"runtime"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func readHomeTimeline(ctx context.Context, req *social.ReadHomeTimelineRequest) *social.ReadHomeTimelineResponse {
	slowpoke.SlowpokeCheck("readHomeTimeline")
	posts := social.ReadHomeTimeline(ctx, req.UserId)
	//fmt.Printf("Posts read: %+v\n", posts)
	resp := social.ReadHomeTimelineResponse{Posts: posts}
	return &resp
}

func writeHomeTimeline(ctx context.Context, req *social.WriteHomeTimelineRequest) *string {
	slowpoke.SlowpokeCheck("writeHomeTimeline")
	social.WriteHomeTimeline(ctx, req.UserId, req.PostIds)
	resp := "OK"
	return &resp
}

func main() {
	if common.ShardEnabled {
		fmt.Println(runtime.GOMAXPROCS(1))
	} else {
		fmt.Println(runtime.GOMAXPROCS(8))
	}
	fmt.Println("Max procs: ", runtime.GOMAXPROCS(0))
	// go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_read_home_timeline", wrappers.ROWrapper[social.ReadHomeTimelineRequest, social.ReadHomeTimelineResponse](readHomeTimeline))
	http.HandleFunc("/write_home_timeline", wrappers.NonROWrapper[social.WriteHomeTimelineRequest, string](writeHomeTimeline))
	slowpoke.SlowpokeInit()
	fmt.Println("Starting server on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
