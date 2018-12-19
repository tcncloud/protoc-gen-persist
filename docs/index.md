## protoc-gen-persist

protoc-gen-persist is a protoc plugin.
Its goal is to help write your persistence layer for simple GRPC calls that mostly deal with
talking to a database.

protoc-gen-persist will look at a protobuf service's options and generate implementation for
the service's methods.

you activate the plugin using protoc's CLI.
```bash
protoc -I. -I$GOPATH/src --persist_out=plugins=protoc-gen-persist:$GOPATH/src ./*.proto
```

if you do not have fully qualified go_package options in your proto files while running protoc
from your $GOPATH/src,  you will need to specify the persist_root option. Which is the
go package base that the persist_lib will extend from.

The generated persist file needs to import persist_lib, so specify the root of the package
to do that.
```bash
protoc -I. -I$GOPATH/src --persist_out=plugins=protoc-gen-persist,persist_root=github.com/protoc-gen-persist/examples/user_sql/pb:. ./pb/*.proto
```

the persist plugin will generate service handlers that implement a grpc service.
It looks at a method's stream type, input message, output message, and a few user
specified options, and decides how the database must be structured.


It then writes function to marshal a protobuf message to, and from the database row,
perform iterations over a protobuf from a database's iterator, and functions that
run the protobuf option's query on the backend.
for example, given this [snippet from our sql examples](https://github.com/tcncloud/protoc-gen-persist/blob/master/examples/user_sql/main.go)
```proto
message Friends {
	repeated string names = 1;
}

message User {
	int64 id = 1;
	string name = 2;
	Friends friends = 3;
	google.protobuf.Timestamp created_on = 4;
}

service UServ {
	option (persist.service_type) = SQL;
	rpc SelectUserById(User) returns (User) {
		option (persist.ql) = {
			query: ["SELECT id, name, friends, created_on FROM users WHERE id = $1"],
			arguments: ["id"],
		};
	};
}
```

you get the following go code:

in [the pb package](https://github.com/tcncloud/protoc-gen-persist/tree/master/examples/user_sql/pb):
- function that marshals from protobuf message to a database row
```go
func UserToUServPersistType(serv UServTypeMapping, req *User) (*persist_lib.UserForUServ, error) {
	params := &persist_lib.UserForUServ{}
	params.Id = req.Id
	params.Name = req.Name
	if req.Friends == nil {
		req.Friends = new(Friends)
	}
	{
		raw, err := proto.Marshal(req.Friends)
		if err != nil {
			return nil, err
		}
		params.Friends = raw
	}
	mapper := serv.TimestampTimestamp()
	params.CreatedOn = mapper.ToSql(req.CreatedOn)
	return params, nil
}
```
- function that marshals from database row to protobuf message
```go
func UserFromUServDatabaseRow(serv UServTypeMapping, row persist_lib.Scanable) (*User, error) {
	res := &User{}
	var Id_ int64
	var Name_ string
	var Friends_ []byte
	CreatedOn_ := serv.TimestampTimestamp().Empty()
	if err := row.Scan(
		&Id_,
		&Name_,
		&Friends_,
		CreatedOn_,
	); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	res.Id = Id_
	res.Name = Name_
	{
		var converted = new(Friends)
		if err := proto.Unmarshal(Friends_, converted); err != nil {
			return nil, err
		}
		res.Friends = converted
	}
	if err := CreatedOn_.ToProto(&res.CreatedOn); err != nil {
		return nil, err
	}
	return res, nil
}
```
- function that iterates over a database iterator returning protobuf
messages
```go
func IterUServUserProto(ms UServTypeMapping, iter *persist_lib.Result, next func(i *User) error) error {
	return iter.Do(func(r persist_lib.Scanable) error {
		item, err := UserFromUServDatabaseRow(ms, r)
		if err != nil {
			return fmt.Errorf("error converting User row to protobuf message: %s", err)
		}
		return next(item)
	})
}

```
- struct implements the service for the protobuf method
```go
type UServImpl struct {
	PERSIST   *persist_lib.UServMethodReceiver
	FORWARDED RestOfUServHandlers
	HOOKS     UServHooks
	MAPPINGS  UServTypeMapping
}
// ...
func (s *UServImpl) SelectUserById(ctx context.Context, req *User) (*User, error) {
	var err error
	var res = &User{}
	_ = err
	_ = res
	params, err := UserToUServPersistType(req)
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.SelectUserById(ctx, params, func(row persist_lib.Scanable) {
		if row == nil { // there was no return data
			return
		}
		res, err = UserFromUServDatabaseRow(row)
		if err != nil {
			iterErr = err
			return
		}
	})
	if err != nil {
		return nil, gstatus.Errorf(codes.Unknown, "error calling persist service: %v", err)
	} else if iterErr != nil {
		return nil, iterErr
	}
	return res, nil
}
```
- builder for the  struct that implements the service
```go

type UServImplBuilder struct {
	err           error
	rest          RestOfUServHandlers
	queryHandlers *persist_lib.UServQueryHandlers
	i             *UServImpl
	db            sql.DB
	hooks         UServHooks
	mappings      UServTypeMapping
}

func NewUServBuilder() *UServImplBuilder {
	return &UServImplBuilder{i: &UServImpl{}}
}
func (b *UServImplBuilder) WithHooks(hs UServHooks) *UServImplBuilder {
	b.hooks = hs
	return b
}
func (b *UServImplBuilder) WithTypeMapping(ts UServTypeMapping) *UServImplBuilder {
	b.mappings = ts
	return b
}
func (b *UServImplBuilder) WithRestOfGrpcHandlers(r RestOfUServHandlers) *UServImplBuilder {
	b.rest = r
	return b
}
func (b *UServImplBuilder) WithPersistQueryHandlers(p *persist_lib.UServQueryHandlers) *UServImplBuilder {
	b.queryHandlers = p
	return b
}
func (b *UServImplBuilder) WithDefaultQueryHandlers() *UServImplBuilder {
	accessor := persist_lib.NewSqlClientGetter(&b.db)
	queryHandlers := &persist_lib.UServQueryHandlers{
		SelectUserByIdHandler:  persist_lib.DefaultSelectUserByIdHandler(accessor),
	}
	b.queryHandlers = queryHandlers
	return b
}

// set the custom handlers you want to use in the handlers
// this method will make sure to use a default handler if
// the handler is nil.
func (b *UServImplBuilder) WithNilAsDefaultQueryHandlers(p *persist_lib.UServQueryHandlers) *UServImplBuilder {
	accessor := persist_lib.NewSqlClientGetter(&b.db)
	if p.SelectUserByIdHandler == nil {
		p.SelectUserByIdHandler = persist_lib.DefaultSelectUserByIdHandler(accessor)
	}
	b.queryHandlers = p
	return b
}
func (b *UServImplBuilder) WithSqlClient(c *sql.DB) *UServImplBuilder {
	b.db = *c
	return b
}
func (b *UServImplBuilder) WithNewSqlDb(driverName, dataSourceName string) *UServImplBuilder {
	db, err := sql.Open(driverName, dataSourceName)
	b.err = err
	b.db = *db
	return b
}
func (b *UServImplBuilder) Build() (*UServImpl, error) {
	if b.err != nil {
		return nil, b.err
	}
	b.i.PERSIST = &persist_lib.UServMethodReceiver{Handlers: *b.queryHandlers}
	b.i.FORWARDED = b.rest
	b.i.HOOKS = b.hooks
	b.i.MAPPINGS = b.mappings
	return b.i, nil
}
func (b *UServImplBuilder) MustBuild() *UServImpl {
	s, err := b.Build()
	if err != nil {
		panic("error in builder: " + err.Error())
	}
	return s
}
```

in [the persist_lib package](https://github.com/tcncloud/protoc-gen-persist/tree/master/examples/user_sql/pb/persist_lib):
-  struct that matches protobuf type with getters and setters
```go
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
```
- interface for the query consisting of getters for the query parameters
```go
type UServSelectUserByIdQueryParams interface {
	GetId() int64
}
```
- function that performs the database query on the iterface matched by the
protobuf method's mapped input message type
```go
func UServSelectUserByIdQuery(tx Runable, req UServSelectUserByIdQueryParams) *Result {
	row := tx.QueryRow(
		"SELECT id, name, friends, created_on FROM users WHERE id = $1 ",
		req.GetId(),
	)
	return newResultFromRow(row)
}
```
- default implemented query handler used by the builder
```go
type UServMethodReceiver struct {
	Handlers UServQueryHandlers
}
type UServQueryHandlers struct {
	SelectUserByIdHandler  func(context.Context, *UserForUServ, func(Scanable)) error
}
// next must be called on each result row
func (p *UServMethodReceiver) SelectUserById(ctx context.Context, params *UserForUServ, next func(Scanable)) error {
	return p.Handlers.SelectUserByIdHandler(ctx, params, next)
}

func DefaultSelectUserByIdHandler(accessor SqlClientGetter) func(context.Context, *UserForUServ, func(Scanable)) error {
	return func(ctx context.Context, req *UserForUServ, next func(Scanable)) error {
		sqlDB, err := accessor()
		if err != nil {
			return err
		}
		res := UServSelectUserByIdQuery(sqlDB, req)
		err = res.Do(func(row Scanable) error {
			next(row)
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	}
}
```

How all this fits together can be understood better by looking at the examples


protoc-gen-persist can generate code for sql, and spanner backends.

- for sql it uses go's [database/sql package](https://golang.org/pkg/database/sql/).
It is not driver specific.
- for spanner it uses [google's golang spanner sdk](https://godoc.org/cloud.google.com/go/spanner)


Persist is opinionated.  It decides how to fit a protobuf message into the database, based on
the types of the fields in the request/response message.  If a type does not fit in the
database, it either converts the code to something it knows will fit,  or generates
incorrect code.


If the default mappings for protobuf type do not work for your database row type, you
can provide a custom mapping, but filling out a our type mapping option on the service.

### default mappings for types are


#### for spanner
- enums are transformed into an int64
- message types are transformed into []byte
- repeated message types are transformed into [][]byte
- repeated string is transformed to []spanner.NullString and finally into []string
- repeated int64 is transformed to []spanner.NullInt64 and finally into []int64
- repeated bool  is transformed to []spanner.NullBool and finally into []bool
- repeated float64 is transformed into []spanner.NullFloat64[] and finally into float64
- float64, string, and int64 types are left unconverted (because they fit)

repeated enums are not supported, and will require a custom type mapping
all other types are not supported and will require a custom type mapping


#### for sql
- enums are transformed into an int32
- message types are stored as []byte
- int32, int64, bool, float32, float64, and string are left unconverted (because they fit)


any repeated type will not be satisfied automatically by the driver.Value interface
so it does not have a supported default mapping.
(Even if your database driver can fit that type)


driver.Value types are:
- int64
- float64
- bool
- []byte
- string
- time.Time


custom type mappings can be defined on a service to map encountered protobuf fields to
database types and back.

an example of one:
```proto
service UServ {
	option (persist.service_type) = SQL;
	option (persist.mapping) = {
		types: [{
			proto_type_name: "google.protobuf.Timestamp",
			proto_type: TYPE_MESSAGE,
		},{ proto_type_name: "pb.SliceStringParam",
			proto_type: TYPE_MESSAGE,
		}]
	};
}
```

the protobuf description of the type mapping option:
```proto
message TypeMapping {
    message TypeDescriptor {
        // if this is not setup the proto_type must be one of the built-in types
        optional string proto_type_name =1;
        // If proto_type_name is set, this need not be set.  If both this and proto_type_name
        // are set, this must be one of TYPE_ENUM, TYPE_MESSAGE
        // TYPE_GROUP is not supported
        optional google.protobuf.FieldDescriptorProto.Type proto_type= 2;
        // if proto_label is not setup we consider any option except LABAEL_REPEATED
        optional google.protobuf.FieldDescriptorProto.Label proto_label = 3;
    }
    repeated TypeDescriptor types = 1;
}
```

the code that is generated and needs to be implemented
```go
type UServTypeMapping interface {
	TimestampTimestamp() UServTimestampTimestampMappingImpl
	SliceStringParam() UServSliceStringParamMappingImpl
}
type UServTimestampTimestampMappingImpl interface {
	ToProto(**timestamp.Timestamp) error
	ToSql(*timestamp.Timestamp) sql.Scanner
	Empty() UServTimestampTimestampMappingImpl
	sql.Scanner
	driver.Valuer
}
type UServSliceStringParamMappingImpl interface {
	ToProto(**SliceStringParam) error
	ToSql(*SliceStringParam) sql.Scanner
	Empty() UServSliceStringParamMappingImpl
	sql.Scanner
	driver.Valuer
}
```

to map a type from protobuf to the database, you need to implement a type with 4 methods
a type [in our examples](https://github.com/tcncloud/protoc-gen-persist/blob/master/examples/user_sql/pb/time_converter.go)
for converting google protobuf typestamps looks like this:
```go
type TimeString struct {
	t *timestamp.Timestamp
}
```
- ToSql/ToSpanner  intializes  our type from a protobuf message
```go
func (ts TimeString) ToSql(t *timestamp.Timestamp) *TimeString {
	ts.t = t
	return &ts
}
```
- ToProto must set a stored variable outside of the function 
```go
func (ts TimeString) ToProto(req **timestamp.Timestamp) error {
	*req = ts.t
	return nil
}
```
- (Scan, Value) (if using a SQL backend)  Scan will populate our type from the database
value will convert our type into something that fits in the driver's backend
```go
func (t *TimeString) Scan(src interface{}) error {
	tStr, ok := src.(string)
	if !ok {
		return fmt.Errorf("cannot scan out timestamp from not a string")
	}
	ti, err := time.Parse(time.RFC3339, tStr)
	if err != nil {
		return err
	}
	stamp, err := ptypes.TimestampProto(ti)
	if err != nil {
		return err
	}
	t.t = stamp
	return nil
}

func (t *TimeString) Value() (driver.Value, error) {
	return ptypes.TimestampString(t.t), nil
}
```
- Empty must return a new, empty value of the underlying type that implements the type mapping interface
```go
func (t TimeString) Empty() UServTimestampTimestampMappingImpl {
	return new(TimeString)
}
```
- (SpannerScan, SpannerValue)  (if using a SPANNER backend)  SpannerScan will need to
convert from a spanner.GenericColumnValue, and SpannerValue will need to convert the message
into a value that fits in spanner
```go
func (t *TimeString) SpannerScan(src *spanner.GenericColumnValue) error {
	var tStr string
	if err := src.Decode(&tStr); err != nil {
		return err
	}
	ti, err := time.Parse(time.RFC3339, tStr)
	if err != nil {
		return err
	}
	stamp, err := ptypes.TimestampProto(ti)
	if err != nil {
		return err
	}
	t.t = stamp
	return nil
}

func (t *TimeString) SpannerValue() (interface{}, error) {
	return ptypes.TimestampString(t.t), nil
}
```


### before/after hooks
you can specify optional functions that will hook into your generated service handler
for a method.  One runs before the query  (the before hook)  another runs after
the query, and before results are sent back to the user.  (the after hook)

before hooks need to take the protbuf message used as the input for the method as a parameter
and return the (protobuf message output, error) as a response. (unless the method
is server streaming,  then it needs to return ([]protobuf message output, error))

- if a before hook returns a non nil result, and a nil error, the query will NOT be ran,
and the results will be returned to the client.  This is most useful for caching.
- a pointer to the actual request is given to the before hook.  So any changes made
to the request object will persist over to the handler.  This could be useful for a
number of reasons.  (auto incrementing ids stored on the server for example)
```go
// example hooks where *pb.Name is a input to the protobuf services method,
// and pb.ExampleTable is the output message for the method

// server streaming before hook looks like this
func ServerStreamBeforeHook(req *pb.Name) ([]*pb.ExampleTable, error) {
	fmt.Printf("ServerStreamBeforeHook: %+v\n", req)
	return nil, nil
}

// every other before hook looks like this
func GenericBeforeHook(req *pb.Name) (*pb.ExampleTable, error) {
	fmt.Printf("GenericBeforeHook: %+v\n", req)
	return nil, nil
}
```

after hooks need to take the protobuf message, and protobuf response as parameters and
return error as a response


after hooks are ran on each row recieved from the database. Unlike before hooks,
after hooks are not guaranteed to get the pointer to the request or response message used
in the method.  After hooks main purpose is for their side effects.
```go
// an after hook for the above example's method would look like this
func GenericAfterHook(req *pb.Name, res *pb.ExampleTable) error {
	fmt.Printf("GenericAfterHook: %+v\n", res)
	return nil
}
```
example of protobuf hooks on a method:
```proto
service Test {
  rpc UniaryInsertWithHooks(test.ExampleTable) returns (test.ExampleTable) {
    option (persist.ql) = {
      query: ["insert into example_table (id, start_time, name)  Values (@id, @start_time, \"bananas\")"]
      arguments: ["id", "start_time"]
      before: true
      after: true
    };
  };
}
```

the interface that is generated and needs to be implemented
```go
type UServHooks interface {
	InsertUsersBeforeHook(*User) (*Empty, error)
	InsertUsersAfterHook(*User, *Empty) error
	GetAllUsersBeforeHook(*Empty) ([]*User, error)
	GetAllUsersAfterHook(*Empty, *User) error
}
```


most other questions can be answered by looking at our
[options](https://github.com/tcncloud/protoc-gen-persist/blob/master/persist/options.proto)
or looking at the examples
