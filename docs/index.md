## protoc-gen-persist
A protoc plugin for a persistance layer backing your SQL, or google Spanner database.


protoc-gen-persist (persist) generates go code to help handle serializing, transfering, and iteration over  your protobuf messages that are stored in your database.


## Why persist
When writing our golang grpc services we frequently found that we were writing the same code over and over.
Frequently we had a request protobuf message that we wanted to to send to a service over grpc. The service would:
- receive request
- use request protobuf's fields as input to a query
- perform query
- marshal the resulting rows to the reponse protobuf message's fields
- send response protobuf(s) back over the wire


This code is often pretty similar to other code doing the same thing, but it is easy to get wrong, and hard to organize.

ORMs tie too tightly an object to the database. Refactoring the objects down the road leads to problems, we wanted a looser contract. Our requests and responses should represent inputs and outputs to queries, but should not have to map directly to table definitions.

[given this proto](./petshop.proto) as an example:
Running `$ protoc -I. -I../  --persist_out=. --go_out=plugins=grpc:./test ./*.proto` in this directory

would output (in docs/test) [petshop.persist.go](./test/petshop.persist.go), and [petshop.pb.go](./test/petshop.pb.go)

## Features
- protobuf serialization/deserialization functions for use with golang's [database/sql]() package and [google Spanner]()
- generated grpc service handlers for client streaming, server streaming, and unary methods.
- type mapping for fields on the request/response proto
- "hooks" that can be inserted to run before, or after a query.
- a more robust generated iterator for going over response protos.
- opinionated defaults for handling nested messages, and and implementing grpc methods




## The Basics
These are the options needed for using persist with your grpc service.


the full options can be found [here](../options/options.proto)

```protobuf
extend google.protobuf.ServiceOptions {
    optional QueryOpts ql = 560000;
    optional TypeMapping mapping = 560001;
    optional PersistenceOptions service_type = 560002;
}

enum PersistenceOptions {
    // SQL Query
    SQL = 0;
    SPANNER = 1;
}

message QueryOpts {
    repeated QLImpl queries = 1;
}

message QLImpl {
    // the query string with numeric placeholders for parameters
    // its an array to allow the query to span across multiple lines but it
    // will be joined and used as a single sql query string at generation time
    repeated string query = 1;

    // if provided, persist will rewrite the query string in the generated code
    // replacing "@field_name" (no quotes) with "?" or "$<position>"
    // if unprovided, persist will not rewrite the query string
    optional string pm_strategy = 2;

    // name of this query.  must be unique to the service.
    required string name = 3;

    // the message type that matches the parameters
    // Input rpc messages will be converted to this type
    // they will be used in the parameters in the query
    // The INTERFACE of this message will be used for parameters
    // in the generated query function.
    // if absent, this query takes no  parameters.
    // The query does not have to use all the fields of this type as parameters,
    // but it cannot use any parameter NOT listed here.
    optional string in = 4;

    // the message type that matches what the query returns.
    // This entity message will be converted to the output type
    // input/output messages on rpc calls will have their fields ignored if they
    // don't match this entity.
    // the generated query function will return this message type
    // if absent, this query returns nothing.
    // The query does not have to return a fully populated message,
    // but additional rows returned from the query that do NOT exist on
    // the out message will be ignored.
    optional string out= 5;
}
```

The above code is what is needed to generate the code for persisting data into a database using the `in`put  proto, and iterating over `out`put protos.  **This section of generated code requires no grpc service or methods**.

implementing this option looks like this:
```protobuf
service PetShop{
    option (persist.service_type) = SQL;
    option (persist.ql) = {
        queries: [
            {
                name: "GetCatByName",
                query: [
                    "SELECT",
                        "name,",
                        "age,",
                        "cost",
                    "FROM cats",
                    "WHERE",
                        "name = @cat_name"
                ],
                pm_strategy: "$",
                in: ".test.CatName",
                out: ".test.Cat",
            },
            {
                name: "InsertFish",
                query: [
                    "INSERT INTO fish(",
                        "species,",
                        "cost",
                    ") VALUES(",
                        "@species,",
                        "@cost",
                    ")"
                ],
                pm_strategy: "$",
                in: ".test.Fish",
                out: ".test.Empty",
            }
        ];
    };
}
```

when calling `protoc` in your petshop.persist.go file you will have a struct titled `Queries_Petshop`
__(Queries_SERVICE_NAME)__, with a methods `GetCatByName` and `InsertFish` __(one method for each query name)__

you gain access to this struct by calling the function with this signature:
`func QueriesPetShop(opts ...Opts_PetShop) *Queries_PetShop`


The `opts` here refers to the struct that contains your type mapping initializer, and hooks interfaces.
They are **only** required to be passed in if you are using hooks, or type mapping, and then, they **must** be passed.

The opts struct for our PetShop service has exported fields, so can be initailized directly.  You can also use the generated function. 
```go
type Opts_PetShop struct {
	MAPPINGS TypeMappings_PetShop
	HOOKS    Hooks_PetShop
}

func OptsPetShop(hooks Hooks_PetShop, mappings TypeMappings_PetShop) Opts_PetShop
```

Using these will let you interact with your database using protobufs


Here is an example query, and an of how to use it. 

```protobuf
service PetShop{
    option (persist.service_type) = SQL;
    option(persist.ql) = { 
        queries: [
            {
                name: "GetCatByName",
                query: [
                    "SELECT",
                        "name,",
                        "age,",
                        "cost",
                    "FROM cats",
                    "WHERE",
                        "name = @cat_name"
                ],
                pm_strategy: "$",
                in: ".test.CatName",
                out: ".test.Cat",
            }
        ]
    };
}
```

## Query, Iterator and Rows

This main file shows how to perform a basic query, iterate over the results, and unmarshal it to the output proto.
```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pb "github.com/tcncloud/protoc-gen-persist/docs/test"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	// gain access to the struct that can run this query.
	// 'db' can be anything that satisfies this interface
	/*
		type Runnable interface {
			QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
			ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
		}
	*/
	catQuery := pb.QueriesPetShop().GetCatByName(ctx, db)

	// look for the cat that matches this protobuf's cat_name field
	catName := &pb.CatName{CatName: "boomi"}

	// 'Execute' will always return an iterator.
	iter := catQuery.Execute(catName)

	// you can loop through an iterator calling `Next()`
	// ok will be true if there is something next can return besides EOF
	for {
		row, ok := iter.Next()
		if !ok {
			break
		}
		// you can get the protobuf, and any defered error by calling 'Proto'
		cat, err := row.Proto()

		fmt.Printf("#+v\n", cat)
	}

	// you can loop through an iterator calling Each()
	err = iter.Each(func(row *pb.Row_PetShop_GetCatByName) error {
		// if you have an initialized proto to use already, you can just call unwrap
		cat := pb.Cat{}
		if err := row.Unwrap(&cat); err != nil {
			return err
		}
		// use 'cat'
		fmt.Printf("#+v\n", cat)

		return nil
	})

	// you can assert that an iterator has no results or errors with Zero(), which will error out if there
	// is more than zero rows, or there is an error
	err = iter.Zero()

	// you can assert that an iterator has exactly one result with One(), which will always return one row.
	// The row will always be an error if there is more or less than one result
	row := iter.One()
	// each row will also have the name of the output proto as a dup method of Proto() for clarity

	cat, err := row.Cat()

}
```

### Generated GRPC Handlers
Persist has options that will let it generate a grpc handler for you. Just name the query you want the handler to execute.
```protobuf
service PetShop{
    option (persist.service_type) = SQL;

    // persist will now generate the unary handler implementation that will execute this query
    rpc GetCatByName(CatName) returns (Cat){
        option (persist.opts) = {
            query: "GetCatByName"
        };
    }
    // persist will **ignore** this query, and expect you to write it.
    rpc PetDog(Dog) returns (Empty);

    option(persist.ql) = { 
        queries: [
            {
                name: "GetCatByName",
                query: [
                    "SELECT",
                        "name,",
                        "age,",
                        "cost",
                    "FROM cats",
                    "WHERE",
                        "name = @cat_name"
                ],
                pm_strategy: "$",
                in: ".test.CatName",
                out: ".test.Cat",
            }
        ]
    };
}

```

**RULES**
- Unary handlers take the request as input for the query, and return exactly one output
    - if your query's `out` option has **zero** fields (is empty), then the handler expects no rows
- Server Streaming handlers take the request as input for the query, and stream back the rows as the handler's output proto type.  **Only use generated server streaming handlers if you return rows**
- Client Streaming handlers execute the query on each input streamed to the server
- Bidirectional Streaming is **unsupported**

```protobuf
extend google.protobuf.MethodOptions {
    optional MOpts opts = 560004;
}
message MOpts {
    // must match a name of a QLImpl query in the service.
    required string query = 1;

    // the before function will be called before running any sql code for
    // every input data element and if the return will be not empty/nil and
    // the error will be nil the data returned by this function will be
    // returned by the function skipping the code execution
    optional bool before = 10;

    // the after function will be called after running any sql code for
    // every output data element, the return data of this function will be ignored
    optional bool after = 11;
}
```

**this shows how to start a grpc service with part generated, part custom code.**
```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	pb "github.com/tcncloud/protoc-gen-persist/docs/test"
    "google.golang.org/grpc"
)

func main() {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	// we couldnt have the PetDog Handler generated, so we have to have
	// a struct that implements it as a method
	handlers := &MyPartOfPetShop{}

	// 'ImplPetShop' gets the *FULL* handler implementation of our service handlers.
	service := pb.ImplPetShop(db, handlers)

	// create a grpc server
	server := grpc.NewServer()

	// register our service on our server
	pb.RegisterPetShopServer(server, service)

	
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic(err)
	}
	if err := server.Serve(lis); err != nil {
		fmt.Printf("error serving: %v\n", err)
	}
}
```

## Type Mapping

Persist can *map* fields between types if type_mapping options are implemented. THis is usually done when a field cannot fit naturally in the database. As an example, lets say we have this POSTGRES table, and this proto
we want to store as a row:

**table**
```sql
CREATE TABLE dog_and_fish_owners (
    id TEXT NOT NULL,
    aquarium BYTEA NOT NULL,
    dog_ids TEXT[],
PRIMARY KEY (id));`,
```

**proto**
```protobuf
message Owner {
    string id = 1;
    repeated Cat cats = 2;
    FishTank aquarium = 3;
    DogIds dog_ids = 4;
    double money = 5;
}
service PetShop{
    option(persist.ql) = {
        queries: [
            {
                name: "GetAllOwners",
                query: [ "SELECT id, aquarium, dog_ids FROM dog_and_fish_owners" ],
                pm_strategy: "$",
                in: ".test.Empty",
                out: ".test.Owner",
            }
        ];
    };
}
```

**differences**
- *no cats column.*
    - This will generate not compilable code, because cats is repeated.  Unless working with Spanner, repeated types of any kind are unsupported (yet)
- *aquarium is stored as a byte array*
    - this does **not** need a type mapping.  The default case for persist is to marshal the proto to a `[]byte`
        using proto.Marshal.
- *dog_ids is not a []TEXT*  This case we need a type mapping for. []string is not a type that `database/sql`
    can handle, even though postgres supports it.
- *no money column*
    - this will generate code that will more than likely never get used, but is okay, because unlike the `cats`
        field, `money` is not repeated.


This section will go over how to handle `dog_ids`, which will need a TypeMapping.


Mapping fields from one type to another requires filling out type mapping options.


Type mapping option definition:
```protobuf

extend google.protobuf.ServiceOptions {
    optional TypeMapping mapping = 560001;
}

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

To correctly map your type, you need the full qualified name the protobuf type you are going to map, and mark it either **TYPE_MESSAGE**, or **TYPE_ENUM**.
```protobuf

// example message to map
message DogIds {
    repeated string values = 1;
}

// we have a field that matches one of our mapped types (field #4)
message Owner {
    string id = 1;
    repeated Cat cats = 2;
    FishTank aquarium = 3;
    DogIds dog_ids = 4;
    double money = 5;
}
service PetShop{
    option (persist.mapping) = {
        types: [
            {
                proto_type_name: ".test.DogIds",
                proto_type: TYPE_MESSAGE
            }
        ]
    };
}
```

With our options specified, persist will generate two interfaces we need to implement, and pass as `PetShop_Opts` when using queries containing fields of our mapped type.


**First:**  this interface must return a new instance of our type mapper (implemented as: `DogIder`)
```go
// matches the name of the proto message we are mapping.  If in a different package, the package name is prefixed.
// (matches: proto_type_name: ".test.DogIds",)
type TypeMappings_PetShop interface {
	DogIds() MappingImpl_PetShop_DogIds
}

// our mapping implementation must also implement this interface
type MappingImpl_PetShop_DogIds interface {
	ToProto(**DogIds) error
	ToSql(*DogIds) sql.Scanner
	sql.Scanner
	driver.Valuer
}
```

Here is an example of implementing both needed interfaces:
```go
import (
	"database/sql"
    "database/sql/driver"

	"github.com/lib/pq"
)
// this implements the 'TypeMappings_PetShop` interface
type MyMapper struct{}
func (m *MyMapper) DogIds() pb.MappingImpl_PetShop_DogIds {
    // return a new instance of our `DogIds` Mapper
    return &DogIder{}
}

// This implements the MappingImpl_PetShop_DogIds interface
type DogIder struct{
    // we will need to store this 
    ids *pb.DogIds
}

// ToProto will *only* be called after a call to Scan on the same instance.
// ToProto expects the struct to set its internal parsed proto value from scan into dest (*dest = this.ids)
func (this *DogIder) ToProto(dest **pb.DogIds) error {
    *dest = this.ids
	return nil
}
// ToSql initializes the implementation with a value.  It expects its return value
// to be able to scan out our request as a type that will fit in the db column
func (this *DogIder) ToSql(request *pb.DogIds) sql.Scanner {
    this.ids = request 
	return this
}
// given this arbitrary value, scan needs to populate its internal state so a call to `ToProto` or `Value`
// will return `value` correctly
func (this *DogIder) Scan(value interface{}) error {
    // lib/pq is the driver we use for postgres at TCN. It implments scanner and valuer.
	var in pq.StringArray
	if err := in.Scan(src); err != nil {
		return err
	}
	in.ids = &pb.DogIds{Values: []string(in)}
	return nil
}

// it needs to return its internal value as one that can fit in database/sql driver.Valuer interface
func (this *DogIder) Value() (driver.Value, error) {
    // grab the nested string slice field, hand to pq, return it's call Value()
	return pq.StringArray(this.ids.Values).Value()
}
```

here is an example of it being used
```go
package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"

	_ "github.com/lib/pq"
	pb "github.com/tcncloud/protoc-gen-persist/docs/test"
)
func main() {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	var owners *pb.Owner
	// create an instance of our TypeMapping initializer interface `MyMapper`
	// since we are about to perform a query that uses mapped types, we must create these options
	// to initialize the query with
	opts := pb.OptsPetShop(nil, &MyMapper{})
	// pass opts to the queries constructor.  Anything using this query handler will now use the opts above
	pb.QueriesPetShop(opts).GetAllOwners(ctx, db).Execute().Each(func(row *pb.Row_PetShop_GetAllOwners) error {
		// make an instance of our result
		owner := pb.Owner{}
		// pass to unwrap
		if err := row.Unwrap(&owner); err != nil {
			return err
		}
		// store with rest of owners
		owners = append(owners, owner)

		return nil
	})
	// owners has been successfully mapped
}

```
Spanner will need the spanner equivalent of these functions, which can be found in the spanner examples.



## Before/After hooks
For expensive queries, it might be useful to only run conditionally, and only return cached results from a handler.

- before hooks need to take the protbuf message used as the input for the method as a parameter
and return the (protobuf message output, error) as a response. (unless the method
is server streaming,  then it needs to return ([]protobuf message output, error))
- if a before hook returns a non nil result, and a nil error, the query will NOT be ran,
and the results will be returned to the client.  This is most useful for caching.
- after hooks need to take the protobuf message, and protobuf response as parameters and
return error as a response
- after hooks are ran on each row recieved from the database. Unlike before hooks,
    after hooks are not guaranteed to get the pointer to the request or response message used
    in the method.  After hooks main purpose is for their side effects.









