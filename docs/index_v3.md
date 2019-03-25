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
    // if absent, this query returns nothing, and .
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
	catName := &pb.CatName{CatName: "boomie"}

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





