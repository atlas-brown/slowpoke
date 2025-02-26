package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/hotel"
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

func registerUser(ctx context.Context, req *hotel.RegisterUserRequest) *hotel.RegisterUserResponse {
    slowpoke.SlowpokeCheck("registerUser");
	ok := hotel.RegisterUser(ctx, req.Username, req.Password)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := hotel.RegisterUserResponse{Ok: ok}
	return &resp
}

func login(ctx context.Context, req *hotel.LoginRequest) *hotel.LoginResponse {
    slowpoke.SlowpokeCheck("login");
	token := hotel.Login(ctx, req.Username, req.Password)
	//fmt.Println("Movie info stored for id: " + movieId)
	resp := hotel.LoginResponse{Token: token}
	return &resp
}

func main() {
    slowpoke.SlowpokeCheck("main");
	fmt.Println(runtime.GOMAXPROCS(8))
	slowpoke.SlowpokeInit()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/register_user", wrappers.NonROWrapper[hotel.RegisterUserRequest, hotel.RegisterUserResponse](registerUser))
	http.HandleFunc("/login", wrappers.NonROWrapper[hotel.LoginRequest, hotel.LoginResponse](login))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
