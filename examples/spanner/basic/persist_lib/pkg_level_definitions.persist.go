package persist_lib

import (
	"cloud.google.com/go/spanner"
)

type SpannerClientGetter func() (*spanner.Client, error)

func NewSpannerClientGetter(cli *spanner.Client) SpannerClientGetter {
	return func() (*spanner.Client, error) {
		return cli, nil
	}
}

type Test_NumRowsForExtraSrv struct {
	Count int64
}

// this could be used in a query, so generate the getters/setters
func (p *Test_NumRowsForExtraSrv) GetCount() int64      { return p.Count }
func (p *Test_NumRowsForExtraSrv) SetCount(param int64) { p.Count = param }

type Test_ExampleTableForExtraSrv struct {
	Id        int64
	StartTime []byte
	Name      string
}

// this could be used in a query, so generate the getters/setters
func (p *Test_ExampleTableForExtraSrv) GetId() int64              { return p.Id }
func (p *Test_ExampleTableForExtraSrv) SetId(param int64)         { p.Id = param }
func (p *Test_ExampleTableForExtraSrv) GetStartTime() []byte      { return p.StartTime }
func (p *Test_ExampleTableForExtraSrv) SetStartTime(param []byte) { p.StartTime = param }
func (p *Test_ExampleTableForExtraSrv) GetName() string           { return p.Name }
func (p *Test_ExampleTableForExtraSrv) SetName(param string)      { p.Name = param }

type Test_ExampleTableForMySpanner struct {
	Id        int64
	StartTime interface{}
	Name      string
}

// this could be used in a query, so generate the getters/setters
func (p *Test_ExampleTableForMySpanner) GetId() int64                   { return p.Id }
func (p *Test_ExampleTableForMySpanner) SetId(param int64)              { p.Id = param }
func (p *Test_ExampleTableForMySpanner) GetStartTime() interface{}      { return p.StartTime }
func (p *Test_ExampleTableForMySpanner) SetStartTime(param interface{}) { p.StartTime = param }
func (p *Test_ExampleTableForMySpanner) GetName() string                { return p.Name }
func (p *Test_ExampleTableForMySpanner) SetName(param string)           { p.Name = param }

type SomethingForMySpanner struct {
	Thing []byte
}

// this could be used in a query, so generate the getters/setters
func (p *SomethingForMySpanner) GetThing() []byte      { return p.Thing }
func (p *SomethingForMySpanner) SetThing(param []byte) { p.Thing = param }

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

// this could be used in a query, so generate the getters/setters
func (p *HasTimestampForMySpanner) GetTime() interface{}      { return p.Time }
func (p *HasTimestampForMySpanner) SetTime(param interface{}) { p.Time = param }
func (p *HasTimestampForMySpanner) GetSome() []byte           { return p.Some }
func (p *HasTimestampForMySpanner) SetSome(param []byte)      { p.Some = param }
func (p *HasTimestampForMySpanner) GetStr() string            { return p.Str }
func (p *HasTimestampForMySpanner) SetStr(param string)       { p.Str = param }
func (p *HasTimestampForMySpanner) GetTable() []byte          { return p.Table }
func (p *HasTimestampForMySpanner) SetTable(param []byte)     { p.Table = param }
func (p *HasTimestampForMySpanner) GetStrs() []string         { return p.Strs }
func (p *HasTimestampForMySpanner) SetStrs(param []string)    { p.Strs = param }
func (p *HasTimestampForMySpanner) GetTables() [][]byte       { return p.Tables }
func (p *HasTimestampForMySpanner) SetTables(param [][]byte)  { p.Tables = param }
func (p *HasTimestampForMySpanner) GetSomes() [][]byte        { return p.Somes }
func (p *HasTimestampForMySpanner) SetSomes(param [][]byte)   { p.Somes = param }
func (p *HasTimestampForMySpanner) GetTimes() [][]byte        { return p.Times }
func (p *HasTimestampForMySpanner) SetTimes(param [][]byte)   { p.Times = param }

type Test_ExampleTableRangeForMySpanner struct {
	StartId int64
	EndId   int64
}

// this could be used in a query, so generate the getters/setters
func (p *Test_ExampleTableRangeForMySpanner) GetStartId() int64      { return p.StartId }
func (p *Test_ExampleTableRangeForMySpanner) SetStartId(param int64) { p.StartId = param }
func (p *Test_ExampleTableRangeForMySpanner) GetEndId() int64        { return p.EndId }
func (p *Test_ExampleTableRangeForMySpanner) SetEndId(param int64)   { p.EndId = param }

type Test_NameForMySpanner struct {
	Name string
}

// this could be used in a query, so generate the getters/setters
func (p *Test_NameForMySpanner) GetName() string      { return p.Name }
func (p *Test_NameForMySpanner) SetName(param string) { p.Name = param }
