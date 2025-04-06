package main

import (
	"context"
	"fmt"
	"github.com/eniac/mucache/internal/movie"
	"github.com/eniac/mucache/pkg/slowpoke"
	"github.com/eniac/mucache/pkg/wrappers"
	"net/http"
	"runtime"
	"os"
	"github.com/goccy/go-json"
	"strconv"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func storeCastInfo(ctx context.Context, req *movie.StoreCastInfoRequest) *movie.StoreCastInfoResponse {
    // slowpoke.SlowpokeCheck("storeCastInfo");
	movieId := movie.StoreCastInfo(ctx, req.CastId, req.Name, req.Info)
	//fmt.Println("Movie info stored for id: " + movieId)
	resp := movie.StoreCastInfoResponse{CastId: movieId}
	return &resp
}

func readCastInfos(ctx context.Context, req *movie.ReadCastInfosRequest) *movie.ReadCastInfosResponse {
    // slowpoke.SlowpokeCheck("readCastInfos");
	castInfos := movie.ReadCastInfos(ctx, req.CastIds)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.ReadCastInfosResponse{Infos: castInfos}
	return &resp
}

func populate() {
	ctx := context.Background()
    cast_json, err := os.ReadFile("/app/internal/movie/data/casts_1_500.json")
	if err != nil {
		panic(err)
	}
	var data []map[string]interface{}
	err = json.Unmarshal(cast_json, &data)
	if err != nil {
		panic(err)
	}
	for _, cast := range data {
		castId := strconv.Itoa(int(cast["id"].(float64)))
		name := cast["name"].(string)
		info := cast["biography"].(string)
		movie.StoreCastInfo(ctx, castId, name, info)
	}
	fmt.Println("Populated %d movie ids", len(data))
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	populate()
	slowpoke.SlowpokeInit()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/store_cast_info", wrappers.SlowpokeWrapper[movie.StoreCastInfoRequest, movie.StoreCastInfoResponse](storeCastInfo, "storeCastInfo"))
	http.HandleFunc("/ro_read_cast_infos", wrappers.SlowpokeWrapper[movie.ReadCastInfosRequest, movie.ReadCastInfosResponse](readCastInfos, "readCastInfos"))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
