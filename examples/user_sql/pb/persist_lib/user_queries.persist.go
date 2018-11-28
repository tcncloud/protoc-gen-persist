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
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromRows(res)
}
func UServUpdateUserNamesQuery(tx Runable, req UServUpdateUserNamesQueryParams) *Result {
	res, err := tx.Query(
		"Update users set name = @name WHERE id = @id  RETURNING id, name, friends, created_on",
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromRows(res)
}
func UServUpdateNameToFooQuery(tx Runable, req UServUpdateNameToFooQueryParams) *Result {
	res, err := tx.Exec(
		"Update users set name = 'foo' WHERE id = @id",
	)
	if err != nil {
		return newResultFromErr(err)
	}
	return newResultFromSqlResult(res)
}
func UServGetFriendsQuery(tx Runable, req UServGetFriendsQueryParams) *Result {
	res, err := tx.Query(
		"SELECT id, name, friends, created_on FROM users WHERE name = ANY(@names)",
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
}
type UServGetAllUsersQueryParams interface {
}
type UServSelectUserByIdQueryParams interface {
}
type UServUpdateUserNamesQueryParams interface {
}
type UServUpdateNameToFooQueryParams interface {
}
type UServGetFriendsQueryParams interface {
}
type UServDropTableQueryParams interface {
}
