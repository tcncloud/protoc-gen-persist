package persist_lib

import "cloud.google.com/go/spanner"

type SpannerClientGetter func() (*spanner.Client, error)

func NewSpannerClientGetter(cli *spanner.Client) SpannerClientGetter {
	return func() (*spanner.Client, error) {
		return cli, nil
	}
}

type Test_NumRowsForExtraSrv struct {
	Count int64
}
type Test_ExampleTableForExtraSrv struct {
	Id        int64
	StartTime []byte
	Name      string
}
type Test_ExampleTableForMySpanner struct {
	Id        int64
	StartTime interface{}
	Name      string
}
type SomethingForMySpanner struct {
	Thing []byte
}
type HasTimestampForMySpanner struct {
	Time   interface{}
	Some   []byte
	Str    string
	Table  []byte
	Strs   []string
	Tables [][]byte
	Somes  [][]byte
	Times  [][]byte
}
type Test_ExampleTableRangeForMySpanner struct {
	StartId int64
	EndId   int64
}
type Test_NameForMySpanner struct {
	Name string
}
