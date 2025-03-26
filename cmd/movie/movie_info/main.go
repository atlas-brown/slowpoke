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

func storeMovieInfo(ctx context.Context, req *movie.StoreMovieInfoRequest) *movie.StoreMovieInfoResponse {
    slowpoke.SlowpokeCheck("storeMovieInfo");
	movieId := movie.StoreMovieInfo(ctx, req.MovieId, req.Info, req.CastIds, req.PlotId)
	//fmt.Println("Movie info stored for id: " + movieId)
	resp := movie.StoreMovieInfoResponse{MovieId: movieId}
	return &resp
}

func readMovieInfo(ctx context.Context, req *movie.ReadMovieInfoRequest) *movie.ReadMovieInfoResponse {
    slowpoke.SlowpokeCheck("readMovieInfo");
	movieInfo := movie.ReadMovieInfo(ctx, req.MovieId)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.ReadMovieInfoResponse{Info: movieInfo}
	return &resp
}

func populate() {
	ctx := context.Background()
    movie_json, err := os.ReadFile("/app/internal/movie/data/movies_1_500.json")
	if err != nil {
		panic(err)
	}
	var data []map[string]interface{}
	err = json.Unmarshal(movie_json, &data)
	if err != nil {
		panic(err)
	}

	// cast_json, err := os.ReadFile("/app/internal/movie/data/casts_1_500.json")
	// if err != nil {
	// 	panic(err)
	// }
	// var castData []map[string]interface{}
	// err = json.Unmarshal(cast_json, &castData)
	// if err != nil {
	// 	panic(err)
	// }
	// casts = make(map[string]string)
	// for _, cast := range castData {
	// 	castId := strconv.Itoa(int(cast["id"].(float64)))
	// 	name := cast["name"].(string)

	movie_id := 0
	for _, movie_ := range data {
		movidId := strconv.Itoa(movie_id)
		movieInfo := movie_["title"].(string)
		castIds := make([]string, 0)
		for _, cast := range movie_["cast"].([]interface{}) {
			castId := strconv.Itoa(int(cast.(map[string]interface{})["id"].(float64)))
			castIds = append(castIds, castId)
		}
		plotId := strconv.Itoa(movie_id)
		movie.StoreMovieInfo(ctx, movidId, movieInfo, castIds, plotId)
		movie_id++
	}
	fmt.Println("Populated %d movie infos", movie_id)
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	populate()
	slowpoke.SlowpokeInit()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/store_movie_info", wrappers.NonROWrapper[movie.StoreMovieInfoRequest, movie.StoreMovieInfoResponse](storeMovieInfo))
	http.HandleFunc("/ro_read_movie_info", wrappers.ROWrapper[movie.ReadMovieInfoRequest, movie.ReadMovieInfoResponse](readMovieInfo))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
