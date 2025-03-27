package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/boutique"
	// "github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/wrappers"
	"github.com/eniac/mucache/pkg/slowpoke"
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

func placeOrder(ctx context.Context, req *boutique.PlaceOrderRequest) *boutique.PlaceOrderResponse {
	slowpoke.SlowpokeCheck("placeOrder")
	order := boutique.PlaceOrder(ctx, req.UserId, req.UserCurrency, req.Address, req.Email, req.CreditCard)
	resp := boutique.PlaceOrderResponse{Order: order}
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	// go cm.ZmqProxy()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/place_order", wrappers.NonROWrapper[boutique.PlaceOrderRequest, boutique.PlaceOrderResponse](placeOrder))
	slowpoke.SlowpokeInit()
	fmt.Println("Server started on port 3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	slowpokeListener := &slowpoke.SlowpokeListener{listener}
	panic(http.Serve(slowpokeListener, nil))
}
