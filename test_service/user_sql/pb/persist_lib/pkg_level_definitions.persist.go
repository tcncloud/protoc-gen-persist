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

type FriendsQueryForUServ struct {
	Names interface{}
}

// this could be used in a query, so generate the getters/setters
func (p *FriendsQueryForUServ) GetNames() interface{}      { return p.Names }
func (p *FriendsQueryForUServ) SetNames(param interface{}) { p.Names = param }
