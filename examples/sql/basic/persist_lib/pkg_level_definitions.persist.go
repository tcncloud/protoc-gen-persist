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
type Result struct {
	result sql.Result
	row    *sql.Row
	rows   *sql.Rows
	err    error
}

func newResultFromSqlResult(r sql.Result) *Result {
	return &Result{result: r}
}
func newResultFromRow(r *sql.Row) *Result {
	return &Result{row: r}
}
func newResultFromRows(r *sql.Rows) *Result {
	return &Result{rows: r}
}
func newResultFromErr(err error) *Result {
	return &Result{err: err}
}
func (r *Result) Do(fun func(Scanable) error) error {
	if r.err != nil {
		return r.err
	}
	if r.row != nil {
		if err := fun(r.row); err != nil {
			return err
		}
	}
	if r.rows != nil {
		defer r.rows.Close()
		for r.rows.Next() {
			if err := fun(r.rows); err != nil {
				return err
			}
		}
	}
	return nil
}

// returns sql.ErrNoRows if it did not scan into dest
func (r *Result) Scan(dest ...interface{}) error {
	if r.result != nil {
		return sql.ErrNoRows
	} else if r.row != nil {
		return r.row.Scan(dest...)
	} else if r.rows != nil {
		err := r.rows.Scan(dest...)
		if r.rows.Next() {
			r.rows.Close()
		}
		return err
	}
	return sql.ErrNoRows
}
func (r *Result) Err() error {
	return r.err
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
