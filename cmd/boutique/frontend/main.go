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

func home(ctx context.Context, req *boutique.HomeRequest) *boutique.HomeResponse {
	slowpoke.SlowpokeCheck("home")
	resp := boutique.Home(ctx, *req)
	return &resp
}

func setCurrency(ctx context.Context, req *boutique.FrontendSetCurrencyRequest) *boutique.FrontendSetCurrencyResponse {
	slowpoke.SlowpokeCheck("setCurrency")
	boutique.FrontendSetCurrency(ctx, req.Cur)
	resp := boutique.FrontendSetCurrencyResponse{OK: "OK"}
	return &resp
}

func browseProduct(ctx context.Context, req *boutique.BrowseProductRequest) *boutique.BrowseProductResponse {
	slowpoke.SlowpokeCheck("browseProduct")
	resp := boutique.BrowseProduct(ctx, req.ProductId)
	return &resp
}

func addToCart(ctx context.Context, request *boutique.AddToCartRequest) *boutique.AddToCartResponse {
	slowpoke.SlowpokeCheck("addToCart")
	resp := boutique.AddToCart(ctx, *request)
	return &resp
}

func viewCart(ctx context.Context, request *boutique.ViewCartRequest) *boutique.ViewCartResponse {
	slowpoke.SlowpokeCheck("viewCart")
	resp := boutique.ViewCart(ctx, *request)
	return &resp
}

func checkout(ctx context.Context, request *boutique.CheckoutRequest) *boutique.CheckoutResponse {
	slowpoke.SlowpokeCheck("checkout")
	resp := boutique.Checkout(ctx, *request)
	return &resp
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	// go cm.ZmqProxy() // Disable cache proxy
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/ro_home", wrappers.ROWrapper[boutique.HomeRequest, boutique.HomeResponse](home))
	http.HandleFunc("/set_currency", wrappers.NonROWrapper[boutique.FrontendSetCurrencyRequest, boutique.FrontendSetCurrencyResponse](setCurrency))
	http.HandleFunc("/ro_browse_product", wrappers.ROWrapper[boutique.BrowseProductRequest, boutique.BrowseProductResponse](browseProduct))
	http.HandleFunc("/add_to_cart", wrappers.NonROWrapper[boutique.AddToCartRequest, boutique.AddToCartResponse](addToCart))
	http.HandleFunc("/ro_view_cart", wrappers.ROWrapper[boutique.ViewCartRequest, boutique.ViewCartResponse](viewCart))
	http.HandleFunc("/checkout", wrappers.ROWrapper[boutique.CheckoutRequest, boutique.CheckoutResponse](checkout))
	slowpoke.SlowpokeInit()
	fmt.Println("Server started on :3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	slowpokeListener := &slowpoke.SlowpokeListener{listener}
	panic(http.Serve(slowpokeListener, nil))
}
