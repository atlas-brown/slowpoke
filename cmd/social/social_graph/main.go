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

func InsertUser(ctx context.Context, req *social.InsertUserRequest) *string {
	slowpoke.SlowpokeCheck("InsertUser")
	social.InsertUser(ctx, req.UserId)
	resp := "OK"
	return &resp
}

func GetFollowers(ctx context.Context, req *social.GetFollowersRequest) *social.GetFollowersResponse {
	slowpoke.SlowpokeCheck("GetFollowers")
	followers := social.GetFollowers(ctx, req.UserId)
	resp := social.GetFollowersResponse{
		Followers: followers,
	}
	return &resp
}

func GetFollowees(ctx context.Context, req *social.GetFolloweesRequest) *social.GetFolloweesResponse {
	slowpoke.SlowpokeCheck("GetFollowees")
	followees := social.GetFollowees(ctx, req.UserId)
	resp := social.GetFolloweesResponse{
		Followees: followees,
	}
	return &resp
}

func Follow(ctx context.Context, req *social.FollowRequest) *string {
	slowpoke.SlowpokeCheck("Follow")
	social.Follow(ctx, req.FollowerId, req.FolloweeId)
	resp := "OK"
	return &resp
}

func FollowMulti(ctx context.Context, req *social.FollowManyRequest) *string {
	slowpoke.SlowpokeCheck("FollowMulti")
	social.FollowMulti(ctx, req.UserId, req.FollowerIds, req.FolloweeIds)
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
	http.HandleFunc("/insert_user", wrappers.NonROWrapper[social.InsertUserRequest, string](InsertUser))
	http.HandleFunc("/ro_get_followers", wrappers.ROWrapper[social.GetFollowersRequest, social.GetFollowersResponse](GetFollowers))
	http.HandleFunc("/ro_get_followees", wrappers.ROWrapper[social.GetFolloweesRequest, social.GetFolloweesResponse](GetFollowees))
	http.HandleFunc("/follow", wrappers.NonROWrapper[social.FollowRequest, string](Follow))
	http.HandleFunc("/follow_multi", wrappers.NonROWrapper[social.FollowManyRequest, string](FollowMulti))
	slowpoke.SlowpokeInit()
	fmt.Println("Starting server on port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
