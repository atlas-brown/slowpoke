package movie

import (
	"context"
	"github.com/eniac/mucache/pkg/slowpoke"
)

func WritePlot(ctx context.Context, plotId string, plot string) string {
	slowpoke.SetState(ctx, plotId, plot)
	return plotId
}

func ReadPlot(ctx context.Context, plotId string) string {
	plot, err := slowpoke.GetState[string](ctx, plotId)
	if err != nil {
		panic(err)
	}
	return plot
}
