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
	"strings"
	"strconv"
)

func heartbeat(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Heartbeat\n"))
	if err != nil {
		return
	}
}

func writePlot(ctx context.Context, req *movie.WritePlotRequest) *movie.WritePlotResponse {
    slowpoke.SlowpokeCheck("writePlot");
	plotId := movie.WritePlot(ctx, req.PlotId, req.Plot)
	//fmt.Println("Movie info stored for id: " + movieId)
	resp := movie.WritePlotResponse{PlotId: plotId}
	return &resp
}

func readPlot(ctx context.Context, req *movie.ReadPlotRequest) *movie.ReadPlotResponse {
    slowpoke.SlowpokeCheck("readPlot");
	plot := movie.ReadPlot(ctx, req.PlotId)
	//fmt.Printf("Movie info read: %v\n", movieInfo)
	resp := movie.ReadPlotResponse{Plot: plot}
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
		plot_id := strconv.Itoa(movie_id)
		plot := strings.Repeat(movie_["overview"].(string), 20)
		movie.WritePlot(ctx, plot_id, plot)
		movie_id++
	}
	fmt.Println("Populated %d movie plots", movie_id)
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(8))
	populate()
	slowpoke.SlowpokeInit()
	http.HandleFunc("/heartbeat", heartbeat)
	http.HandleFunc("/write_plot", wrappers.NonROWrapper[movie.WritePlotRequest, movie.WritePlotResponse](writePlot))
	http.HandleFunc("/ro_read_plot", wrappers.ROWrapper[movie.ReadPlotRequest, movie.ReadPlotResponse](readPlot))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}
