package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/boutique"
	// "github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/slowpoke"
	"github.com/eniac/mucache/pkg/wrappers"
	"net"
	"net/http"
	"runtime"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func getQuote(ctx context.Context, req *boutique.GetQuoteRequest) *boutique.GetQuoteResponse {
	// slowpoke.SlowpokeCheck("getQuote")
	quote := boutique.GetQuote(ctx, req.Items)
	resp := boutique.GetQuoteResponse{CostUsd: quote}
	return &resp
}

func shipOrder(ctx context.Context, req *boutique.ShipOrderRequest) *boutique.ShipOrderResponse {
	// slowpoke.SlowpokeCheck("shipOrder")
	id := boutique.ShipOrder(ctx, req.Address, req.Items)
	resp := boutique.ShipOrderResponse{TrackingId: id}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	// go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	// http.HandleFunc("/ro_get_quote", wrappers.ROWrapper[boutique.GetQuoteRequest, boutique.GetQuoteResponse](getQuote))
	http.HandleFunc("/ro_get_quote", wrappers.SlowpokeWrapper[boutique.GetQuoteRequest, boutique.GetQuoteResponse](getQuote, "getQuote"))
	// http.HandleFunc("/ship_order", wrappers.NonROWrapper[boutique.ShipOrderRequest, boutique.ShipOrderResponse](shipOrder))
	http.HandleFunc("/ship_order", wrappers.SlowpokeWrapper[boutique.ShipOrderRequest, boutique.ShipOrderResponse](shipOrder, "shipOrder"))
	slowpoke.SlowpokeInit()
	fmt.Println("Server started on port 3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	slowpokeListener := &slowpoke.SlowpokeListener{listener}
	panic(http.Serve(slowpokeListener, nil))
}
