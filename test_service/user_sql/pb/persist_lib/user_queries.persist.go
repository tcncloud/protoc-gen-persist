package persist_lib

import "database/sql"

func EmptyFromCreateTableQuery(tx Runable, req EmptyFromCreateTableQueryParams) (sql.Result, error) {
	return tx.Exec(
		"CREATE TABLE users(id integer PRIMARY KEY, name VARCHAR(50), friends BYTEA, created_on VARCHAR(50)) ",
	)
}
func UserFromInsertUsersQuery(tx Runable, req UserFromInsertUsersQueryParams) (sql.Result, error) {
	return tx.Exec(
		"INSERT INTO users (id, name, friends, created_on) VALUES ($1, $2, $3, $4) ",
		req.GetId(),
		req.GetName(),
		req.GetFriends(),
		req.GetCreatedOn(),
	)
}
func EmptyFromGetAllUsersQuery(tx Runable, req EmptyFromGetAllUsersQueryParams) (*sql.Rows, error) {
	return tx.Query(
		"SELECT id, name, friends, created_on FROM users ",
	)
}
func UserFromSelectUserByIdQuery(tx Runable, req UserFromSelectUserByIdQueryParams) *sql.Row {
	return tx.QueryRow(
		"SELECT id, name, friends, created_on FROM users WHERE id = $1 ",
		req.GetId(),
	)
}
func UserFromUpdateUserNamesQuery(tx Runable, req UserFromUpdateUserNamesQueryParams) *sql.Row {
	return tx.QueryRow(
		"Update users set name = $1 WHERE id = $2  RETURNING id, name, friends, created_on ",
		req.GetName(),
		req.GetId(),
	)
}
func FriendsQueryFromGetFriendsQuery(tx Runable, req FriendsQueryFromGetFriendsQueryParams) (*sql.Rows, error) {
	return tx.Query(
		"SELECT id, name, friends, created_on FROM users WHERE name = ANY($1) ",
		req.GetNames(),
	)
}
func EmptyFromDropTableQuery(tx Runable, req EmptyFromDropTableQueryParams) (sql.Result, error) {
	return tx.Exec(
		"drop table users ",
	)
}

type EmptyFromCreateTableQueryParams interface {
}
type UserFromInsertUsersQueryParams interface {
	GetId() int64
	GetName() string
	GetFriends() []byte
	GetCreatedOn() interface{}
}
type EmptyFromGetAllUsersQueryParams interface {
}
type UserFromSelectUserByIdQueryParams interface {
	GetId() int64
}
type UserFromUpdateUserNamesQueryParams interface {
	GetName() string
	GetId() int64
}
type FriendsQueryFromGetFriendsQueryParams interface {
	GetNames() interface{}
}
type EmptyFromDropTableQueryParams interface {
}
