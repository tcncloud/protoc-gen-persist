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

type EmptyForUServ struct {
}

// this could be used in a query, so generate the getters/setters
type UserForUServ struct {
	Id              int64
	Name            string
	Friends         []byte
	CreatedOn       interface{}
	FavoriteNumbers []int64
}

// this could be used in a query, so generate the getters/setters
func (p *UserForUServ) GetId() int64                     { return p.Id }
func (p *UserForUServ) SetId(param int64)                { p.Id = param }
func (p *UserForUServ) GetName() string                  { return p.Name }
func (p *UserForUServ) SetName(param string)             { p.Name = param }
func (p *UserForUServ) GetFriends() []byte               { return p.Friends }
func (p *UserForUServ) SetFriends(param []byte)          { p.Friends = param }
func (p *UserForUServ) GetCreatedOn() interface{}        { return p.CreatedOn }
func (p *UserForUServ) SetCreatedOn(param interface{})   { p.CreatedOn = param }
func (p *UserForUServ) GetFavoriteNumbers() []int64      { return p.FavoriteNumbers }
func (p *UserForUServ) SetFavoriteNumbers(param []int64) { p.FavoriteNumbers = param }

type FriendsForUServ struct {
	Names []string
}

// this could be used in a query, so generate the getters/setters
func (p *FriendsForUServ) GetNames() []string      { return p.Names }
func (p *FriendsForUServ) SetNames(param []string) { p.Names = param }
