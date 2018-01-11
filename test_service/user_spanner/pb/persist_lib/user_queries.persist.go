package persist_lib

import "cloud.google.com/go/spanner"

func UserFromInsertUsersQuery(req UserFromInsertUsersQueryParams) *spanner.Mutation {
	return spanner.InsertMap("users", map[string]interface{}{
		"id":               req.GetId(),
		"name":             req.GetName(),
		"friends":          req.GetFriends(),
		"created_on":       req.GetCreatedOn(),
		"favorite_numbers": req.GetFavoriteNumbers(),
	})
}
func EmptyFromGetAllUsersQuery(req EmptyFromGetAllUsersQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT id, name, friends, created_on, favorite_numbers FROM users",
		Params: map[string]interface{}{},
	}
}
func UserFromSelectUserByIdQuery(req UserFromSelectUserByIdQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT id, name, friends, created_on, favorite_numbers  FROM users WHERE id = @id",
		Params: map[string]interface{}{
			"id": req.GetId(),
		},
	}
}
func UserFromUpdateUserNamesQuery(req UserFromUpdateUserNamesQueryParams) *spanner.Mutation {
	return spanner.UpdateMap("users", map[string]interface{}{
		"name": req.GetName(),
		"id":   req.GetId(),
	})
}
func FriendsFromGetFriendsQuery(req FriendsFromGetFriendsQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT id, name, friends, created_on, favorite_numbers  FROM users WHERE name IN UNNEST(@names)",
		Params: map[string]interface{}{
			"names": req.GetNames(),
		},
	}
}

type UserFromInsertUsersQueryParams interface {
	GetId() int64
	GetName() string
	GetFriends() []byte
	GetCreatedOn() interface{}
	GetFavoriteNumbers() []int64
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
type FriendsFromGetFriendsQueryParams interface {
	GetNames() []string
}
