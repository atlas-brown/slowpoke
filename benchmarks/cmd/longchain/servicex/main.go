package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/trivial"
	"github.com/eniac/mucache/pkg/slowpoke"
	"github.com/eniac/mucache/pkg/wrappers"
	"net"
	"os"
	"net/http"
)

var snum int;
var nextServ string;

func homeHandler(w http.ResponseWriter, r *http.Request) {
	slowpoke.SlowpokeCheck("home");
	var reply string;
	ctx, _ := wrappers.SetupCtxFromHTTPReq(r, false)
	req := trivial.TrivialRequest{Q: "how are things?"}
	for i := 0; i < snum; i++ {
		serv := fmt.Sprintf("service%d", i);
		resp := slowpoke.Invoke[trivial.TrivialResponse](ctx, serv, "ep2", req)
		reply = resp.A
	} 
	fmt.Fprintf(w, reply)
	if f, ok := w.(http.Flusher); ok {
		f.Flush() // Force buffer to flush
	}
	slowpoke.SlowpokeDelay()
}

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	slowpoke.SlowpokeCheck("home");
// 	var reply string;
// 	if (snum > 0) {
// 		ctx, _ := wrappers.SetupCtxFromHTTPReq(r, false)
// 		req := trivial.TrivialRequest{Q: "how are things?"}
// 		resp := slowpoke.Invoke[trivial.TrivialResponse](ctx, nextServ, "ep1", req)
// 		reply = resp.A
// 	} else {
// 		reply = "hoho"
// 	}
// 	slowpoke.Barrier()
// 	fmt.Fprintf(w, reply)
// 	slowpoke.SlowpokeDelay()
// }

func ep1(ctx context.Context, request *trivial.TrivialRequest) *trivial.TrivialResponse {
	slowpoke.SlowpokeCheck("ep1")
	if (snum > 0) {
		req := trivial.TrivialRequest{Q: "how are things?"}
		slowpoke.Invoke[trivial.TrivialResponse](ctx, nextServ, "ep1", req)
	}
	resp := trivial.TrivialResponse{A: "ok"}
	return &resp
}

func ep2(ctx context.Context, request *trivial.TrivialRequest) *trivial.TrivialResponse {
	slowpoke.SlowpokeCheck("ep1")
	resp := trivial.TrivialResponse{A: "ok"}
	return &resp
}

func main() {
	if env, ok := os.LookupEnv("LONGCHAIN_SERVICE_NUM"); ok {
		fmt.Sscanf(env, "%d", &snum)
		fmt.Printf("SLOWPOKE_SERVICE_NUM=%d\n", snum)
		nextServ = fmt.Sprintf("service%d", snum-1);
	} else {
		fmt.Printf("no service number\n")
		return
	}
	http.HandleFunc("/ep1", wrappers.NonROWrapper[trivial.TrivialRequest, trivial.TrivialResponse](ep1, slowpoke.SlowpokeDelay))
	http.HandleFunc("/ep2", wrappers.NonROWrapper[trivial.TrivialRequest, trivial.TrivialResponse](ep2, slowpoke.SlowpokeDelay))
	http.HandleFunc("/home", homeHandler)
	slowpoke.SlowpokeInit()
	fmt.Println("Server started on :3000")
	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic(err)
	}
	slowpokeListener := &slowpoke.SlowpokeListener{listener}
	panic(http.Serve(slowpokeListener, nil))
}