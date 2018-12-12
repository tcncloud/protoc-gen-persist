package persist_lib

func UServCreateTableQuery(tx Runable, req UServCreateTableQueryParams) *Result {
	res, err := tx.Exec(
		"CREATE TABLE users(id integer PRIMARY KEY, name VARCHAR(50), friends BYTEA, created_on VARCHAR(50))",
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromSqlResult(res)
}
func UServInsertUsersQuery(tx Runable, req UServInsertUsersQueryParams) *Result {
	res, err := tx.Exec(
		"INSERT INTO users (id, name, friends, created_on) VALUES (@id, @name, @friends, @created_on)",
		req.GetId(),
		req.GetName(),
		req.GetFriends(),
		req.GetCreatedOn(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromSqlResult(res)
}
func UServGetAllUsersQuery(tx Runable, req UServGetAllUsersQueryParams) *Result {
	res, err := tx.Query(
		"SELECT id, name, friends, created_on FROM users",
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromRows(res)
}
func UServSelectUserByIdQuery(tx Runable, req UServSelectUserByIdQueryParams) *Result {
	res, err := tx.Query(
		"SELECT id, name, friends, created_on FROM users WHERE id = @id",
		req.GetId(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromRows(res)
}
func UServUpdateUserNamesQuery(tx Runable, req UServUpdateUserNamesQueryParams) *Result {
	res, err := tx.Query(
		"Update users set name = @name PK(id = @id) ",
		req.GetName(),
		req.GetId(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromRows(res)
}
func UServUpdateNameToFooQuery(tx Runable, req UServUpdateNameToFooQueryParams) *Result {
	res, err := tx.Exec(
		"Update users set name = 'foo' PRIMARY_KEY(id = @id)",
		req.GetId(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromSqlResult(res)
}
func UServGetFriendsQuery(tx Runable, req UServGetFriendsQueryParams) *Result {
	res, err := tx.Query(
		"SELECT id, name, friends, created_on  FROM users WHERE name IN UNNEST(@names)",
		req.GetNames(),
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromRows(res)
}
func UServDropTableQuery(tx Runable, req UServDropTableQueryParams) *Result {
	res, err := tx.Exec(
		"drop table users",
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromSqlResult(res)
}

type UServCreateTableQueryParams interface {
}
type UServInsertUsersQueryParams interface {
	GetId() int64
	GetName() string
	GetFriends() []byte
	GetCreatedOn() interface{}
}
type UServGetAllUsersQueryParams interface {
}
type UServSelectUserByIdQueryParams interface {
	GetId() int64
}
type UServUpdateUserNamesQueryParams interface {
	GetName() string
	GetId() int64
}
type UServUpdateNameToFooQueryParams interface {
	GetId() int64
}
type UServGetFriendsQueryParams interface {
	GetNames() interface{}
}
type UServDropTableQueryParams interface {
}
