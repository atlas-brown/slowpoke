package movie

import (
	"context"
	"github.com/eniac/mucache/pkg/slowpoke"
)

func ReadPage(ctx context.Context, movieId string) Page {
	req1 := ReadMovieInfoRequest{MovieId: movieId}
	//fmt.Printf("[Page] Movie id asked: %v\n", movieId)
	movieInfoRes := slowpoke.Invoke[ReadMovieInfoResponse](ctx, "movieinfo", "ro_read_movie_info", req1)
	movieInfo := movieInfoRes.Info

	// TODO: Make them async
	req2 := ReadCastInfosRequest{CastIds: movieInfo.CastIds}
	castInfosRes := slowpoke.Invoke[ReadCastInfosResponse](ctx, "castinfo", "ro_read_cast_infos", req2)
	req3 := ReadPlotRequest{PlotId: movieInfo.PlotId}
	plotRes := slowpoke.Invoke[ReadPlotResponse](ctx, "plot", "ro_read_plot", req3)
	req4 := ReadMovieReviewsRequest{MovieId: movieId}
	reviewsRes := slowpoke.Invoke[ReadMovieReviewsResponse](ctx, "moviereviews", "ro_read_movie_reviews", req4)
	//fmt.Printf("[Page] Reviews read: %v\n", reviewsRes)
	page := Page{
		MovieInfo: movieInfo,
		CastInfos: castInfosRes.Infos,
		Plot:      plotRes.Plot,
		Reviews:   reviewsRes.Reviews,
	}
	return page
}
