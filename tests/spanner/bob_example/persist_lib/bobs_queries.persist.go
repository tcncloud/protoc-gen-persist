package persist_lib

import "cloud.google.com/go/spanner"

func BobsDeleteBobsQuery(req BobsDeleteBobsQueryParams) *spanner.Mutation {
	return spanner.Delete("bob_table", spanner.KeyRange{
		Start: spanner.Key{
			"Bob",
		},
		End: spanner.Key{
			"Bob",
			req.GetStartTime(),
		},
		Kind: spanner.ClosedOpen,
	})
}
func BobsPutBobsQuery(req BobsPutBobsQueryParams) *spanner.Mutation {
	return spanner.InsertMap("bob_table", map[string]interface{}{
		"id":         req.GetId(),
		"name":       req.GetName(),
		"start_time": req.GetStartTime(),
	})
}
func BobsGetBobsQuery(req BobsGetBobsQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL:    "SELECT * from bob_table",
		Params: map[string]interface{}{},
	}
}
func BobsGetPeopleFromNamesQuery(req BobsGetPeopleFromNamesQueryParams) spanner.Statement {
	return spanner.Statement{
		SQL: "SELECT * FROM bob_table WHERE name IN UNNEST(@names)",
		Params: map[string]interface{}{
			"names": req.GetNames(),
		},
	}
}

type BobsDeleteBobsQueryParams interface {
	GetStartTime() interface{}
}
type BobsPutBobsQueryParams interface {
	GetName() string
	GetStartTime() interface{}
	GetId() int64
}
type BobsGetBobsQueryParams interface {
}
type BobsGetPeopleFromNamesQueryParams interface {
	GetNames() []string
}
