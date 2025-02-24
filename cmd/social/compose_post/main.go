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

func ComposePost(ctx context.Context, req *social.ComposePostRequest) *string {
	slowpoke.SlowpokeCheck("ComposePost")
	social.ComposePost(ctx, req.Text, req.CreatorId)
	resp := "OK"
	return &resp
}

func ComposePostMulti(ctx context.Context, req *social.ComposePostMultiRequest) *string {
	slowpoke.SlowpokeCheck("ComposePostMulti")
	social.ComposeMulti(ctx, req.Text, req.Number, req.CreatorId)
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
	http.HandleFunc("/compose_post", wrappers.NonROWrapper[social.ComposePostRequest, string](ComposePost))
	http.HandleFunc("/compose_post_multi", wrappers.NonROWrapper[social.ComposePostMultiRequest, string](ComposePostMulti))
	slowpoke.SlowpokeInit()
	fmt.Println("Starting server on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
