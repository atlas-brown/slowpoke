package main

import (
	// "context"
	"fmt"
	"github.com/eniac/mucache/internal/trivial"
	"github.com/eniac/mucache/pkg/slowpoke"
	"github.com/eniac/mucache/pkg/wrappers"
	"net"
	"net/http"
	"os"
)

var sleep_time int

func homeHandler(w http.ResponseWriter, r *http.Request) {
	slowpoke.SlowpokeCheck("home")
	slowpoke.SlowpokePokerPPDelay()
	ctx, _ := wrappers.SetupCtxFromHTTPReq(r, false)
	slowpoke.CPUSpinTime(sleep_time)
	req := trivial.TrivialRequest{Q: "how are things?"}
	resp := slowpoke.Invoke[trivial.TrivialResponse](ctx, "service2", "ep1", req)
	fmt.Fprintf(w, "Welcome to the Home Page %v!", resp.A)
}

func main() {
	if os.Getenv("SLEEP_TIME") != "" {
		fmt.Sscanf(os.Getenv("SLEEP_TIME"), "%d", &sleep_time)
	}
	fmt.Println("SLEEP_TIME:", sleep_time)
	//http.HandleFunc("/", wrappers.NonROWrapper[struct {}, trivial.TrivialResponse](ep1))
	http.HandleFunc("/", homeHandler)
	slowpoke.SlowpokeInit()
	fmt.Println("Server started on :3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	slowpokeListener := &slowpoke.SlowpokeListener{listener}
	panic(http.Serve(slowpokeListener, nil))
}
