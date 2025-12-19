package metadata

import "context"

type Enricher interface {
	EnrichMovie(ctx context.Context, movieID int64) error
	EnrichSeries(ctx context.Context, seriesID int64) error
}
