package persist_lib

import (
	"cloud.google.com/go/spanner"
)

type SpannerClientGetter func() (*spanner.Client, error)

func NewSpannerClientGetter(cli **spanner.Client) SpannerClientGetter {
	return func() (*spanner.Client, error) {
		return *cli, nil
	}
}

type BobForBobs struct {
	Id        int64
	StartTime interface{}
	Name      string
}

// this could be used in a query, so generate the getters/setters
func (p *BobForBobs) GetId() int64                   { return p.Id }
func (p *BobForBobs) SetId(param int64)              { p.Id = param }
func (p *BobForBobs) GetStartTime() interface{}      { return p.StartTime }
func (p *BobForBobs) SetStartTime(param interface{}) { p.StartTime = param }
func (p *BobForBobs) GetName() string                { return p.Name }
func (p *BobForBobs) SetName(param string)           { p.Name = param }

type EmptyForBobs struct {
}

// this could be used in a query, so generate the getters/setters
type NamesForBobs struct {
	Names []string
}

// this could be used in a query, so generate the getters/setters
func (p *NamesForBobs) GetNames() []string      { return p.Names }
func (p *NamesForBobs) SetNames(param []string) { p.Names = param }
