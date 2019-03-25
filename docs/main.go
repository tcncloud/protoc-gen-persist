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
