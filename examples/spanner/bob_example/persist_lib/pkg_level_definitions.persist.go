package persist_lib

import "cloud.google.com/go/spanner"

type SpannerClientGetter func() (*spanner.Client, error)

func NewSpannerClientGetter(cli *spanner.Client) SpannerClientGetter {
	return func() (*spanner.Client, error) {
		return cli, nil
	}
}

type BobForBobs struct {
	Id        int64
	StartTime interface{}
	Name      string
}
type EmptyForBobs struct {
}
type NamesForBobs struct {
	Names []string
}
