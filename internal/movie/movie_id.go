package movie

import (
	"context"
	"github.com/eniac/mucache/pkg/slowpoke"
)

func RegisterMovieId(ctx context.Context, title string, movieId string) {
	slowpoke.SetState(ctx, title, movieId)
}

func GetMovieId(ctx context.Context, title string) string {
	movieId, err := slowpoke.GetState[string](ctx, title)
	if err != nil {
		panic(err)
	}
	return movieId
}
