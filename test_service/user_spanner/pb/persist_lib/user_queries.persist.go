package persist_lib

import "cloud.google.com/go/spanner"

func UServInsertUsersQuery(req UServInsertUsersQueryParams) *spanner.Mutation {
	return spanner.InsertMap("users", map[string]interface{}{
		"id":               req.GetId(),
		"name":             req.GetName(),
		"friends":          req.GetFriends(),
		"created_on":       req.GetCreatedOn(),
		"favorite_numbers": req.GetFavoriteNumbers(),
	})
}
func UServGetAllUsersQuery(req UServGetAllUsersQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT id, name, friends, created_on, favorite_numbers FROM users",
		Params: map[string]interface{}{},
	}
}
func UServSelectUserByIdQuery(req UServSelectUserByIdQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT id, name, friends, created_on, favorite_numbers  FROM users WHERE id = @id",
		Params: map[string]interface{}{
			"id": req.GetId(),
		},
	}
}
func UServUpdateUserNamesQuery(req UServUpdateUserNamesQueryParams) *spanner.Mutation {
	return spanner.UpdateMap("users", map[string]interface{}{
		"name": req.GetName(),
		"id":   req.GetId(),
	})
}
func UServUpdateNameToFooQuery(req UServUpdateNameToFooQueryParams) *spanner.Mutation {
	return spanner.UpdateMap("users", map[string]interface{}{
		"name": "foo",
		"id":   req.GetId(),
	})
}
func UServGetFriendsQuery(req UServGetFriendsQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT id, name, friends, created_on, favorite_numbers  FROM users WHERE name IN UNNEST(@names)",
		Params: map[string]interface{}{
			"names": req.GetNames(),
		},
	}
}

type UServInsertUsersQueryParams interface {
	GetId() int64
	GetName() string
	GetFriends() []byte
	GetCreatedOn() interface{}
	GetFavoriteNumbers() []int64
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
	GetName() string
}
type UServGetFriendsQueryParams interface {
	GetNames() []string
}
