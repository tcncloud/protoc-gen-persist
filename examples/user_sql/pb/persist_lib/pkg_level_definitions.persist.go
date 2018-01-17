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

type EmptyForUServ struct {
}

// this could be used in a query, so generate the getters/setters
type UserForUServ struct {
	Id        int64
	Name      string
	Friends   []byte
	CreatedOn interface{}
}

// this could be used in a query, so generate the getters/setters
func (p *UserForUServ) GetId() int64                   { return p.Id }
func (p *UserForUServ) SetId(param int64)              { p.Id = param }
func (p *UserForUServ) GetName() string                { return p.Name }
func (p *UserForUServ) SetName(param string)           { p.Name = param }
func (p *UserForUServ) GetFriends() []byte             { return p.Friends }
func (p *UserForUServ) SetFriends(param []byte)        { p.Friends = param }
func (p *UserForUServ) GetCreatedOn() interface{}      { return p.CreatedOn }
func (p *UserForUServ) SetCreatedOn(param interface{}) { p.CreatedOn = param }

type FriendsReqForUServ struct {
	Names interface{}
}

// this could be used in a query, so generate the getters/setters
func (p *FriendsReqForUServ) GetNames() interface{}      { return p.Names }
func (p *FriendsReqForUServ) SetNames(param interface{}) { p.Names = param }
