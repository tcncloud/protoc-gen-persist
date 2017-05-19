package basic

import "cloud.google.com/go/spanner"
type MySpannerImpl struct {
	SpannerDB *spanner.Client
}
