package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/movie"
	"github.com/eniac/mucache/pkg/slowpoke"
	"github.com/eniac/mucache/pkg/wrappers"
	"net/http"
	"runtime"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func readPage(ctx context.Context, req *movie.ReadPageRequest) *movie.ReadPageResponse {
    slowpoke.SlowpokeCheck("readPage");
	page := movie.ReadPage(ctx, req.MovieId)
	//fmt.Printf("Page read: %v\n", page)
	resp := movie.ReadPageResponse{Page: page}
	return &resp
}

func main() {
    slowpoke.SlowpokeCheck("main");
	fmt.Println(runtime.GOMAXPROCS(8))
	slowpoke.SlowpokeInit()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_read_page", wrappers.ROWrapper[movie.ReadPageRequest, movie.ReadPageResponse](readPage))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
