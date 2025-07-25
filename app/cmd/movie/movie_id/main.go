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

func registerMovieId(ctx context.Context, req *movie.RegisterMovieIdRequest) *movie.RegisterMovieIdResponse {
    // slowpoke.SlowpokeCheck("registerMovieId");
	movie.RegisterMovieId(ctx, req.Title, req.MovieId)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.RegisterMovieIdResponse{Ok: "OK"}
	return &resp
}

func getMovieId(ctx context.Context, req *movie.GetMovieIdRequest) *movie.GetMovieIdResponse {
    // slowpoke.SlowpokeCheck("getMovieId");
	movieId := movie.GetMovieId(ctx, req.Title)
	resp := movie.GetMovieIdResponse{MovieId: movieId}
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
	movie_id := 0
	for _, movie_ := range data {
		movidId := strconv.Itoa(movie_id)
		title := movie_["title"].(string)
		movie.RegisterMovieId(ctx, title, movidId)
		movie_id++
	}
	fmt.Println("Populated %d movie ids", movie_id)
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	populate()
	slowpoke.SlowpokeInit()
	http.HandleFunc("/heartbeat", heartbeat)
	// http.HandleFunc("/register_movie_id", wrappers.NonROWrapper[movie.RegisterMovieIdRequest, movie.RegisterMovieIdResponse](registerMovieId))
	http.HandleFunc("/register_movie_id", wrappers.SlowpokeWrapper[movie.RegisterMovieIdRequest, movie.RegisterMovieIdResponse](registerMovieId, "registerMovieId"))
	// http.HandleFunc("/ro_get_movie_id", wrappers.ROWrapper[movie.GetMovieIdRequest, movie.GetMovieIdResponse](getMovieId))
	http.HandleFunc("/ro_get_movie_id", wrappers.SlowpokeWrapper[movie.GetMovieIdRequest, movie.GetMovieIdResponse](getMovieId, "getMovieId"))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
