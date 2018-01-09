package persist_lib

import (
	"database/sql"
)

type SqlClientGetter func() (*sql.DB, error)

func NewSqlClientGetter(cli *sql.DB) SqlClientGetter {
	return func() (*sql.DB, error) {
		return cli, nil
	}
}

type Scanable interface {
	Scan(dest ...interface{}) error
}
type Runable interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Exec(string, ...interface{}) (sql.Result, error)
}
type Test_PartialTableForAmazing struct {
	Id        int64
	StartTime interface{}
}

// this could be used in a query, so generate the getters/setters
func (p *Test_PartialTableForAmazing) GetId() int64                   { return p.Id }
func (p *Test_PartialTableForAmazing) SetId(param int64)              { p.Id = param }
func (p *Test_PartialTableForAmazing) GetStartTime() interface{}      { return p.StartTime }
func (p *Test_PartialTableForAmazing) SetStartTime(param interface{}) { p.StartTime = param }

type Test_NameForAmazing struct {
	Name string
}

// this could be used in a query, so generate the getters/setters
func (p *Test_NameForAmazing) GetName() string      { return p.Name }
func (p *Test_NameForAmazing) SetName(param string) { p.Name = param }

type Test_ExampleTableForAmazing struct {
	Id        int64
	StartTime interface{}
	Name      string
}

// this could be used in a query, so generate the getters/setters
func (p *Test_ExampleTableForAmazing) GetId() int64                   { return p.Id }
func (p *Test_ExampleTableForAmazing) SetId(param int64)              { p.Id = param }
func (p *Test_ExampleTableForAmazing) GetStartTime() interface{}      { return p.StartTime }
func (p *Test_ExampleTableForAmazing) SetStartTime(param interface{}) { p.StartTime = param }
func (p *Test_ExampleTableForAmazing) GetName() string                { return p.Name }
func (p *Test_ExampleTableForAmazing) SetName(param string)           { p.Name = param }
