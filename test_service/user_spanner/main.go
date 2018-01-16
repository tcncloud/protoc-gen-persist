package main

import (
	"cloud.google.com/go/spanner"
	admin "cloud.google.com/go/spanner/admin/database/apiv1"
	"encoding/json"
	"fmt"
	"github.com/tcncloud/protoc-gen-persist/test_service/user_spanner/pb"
	pl "github.com/tcncloud/protoc-gen-persist/test_service/user_spanner/pb/persist_lib"
	"golang.org/x/net/context"
	db "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/grpc"
	"io/ioutil"
	"net"
	"os/user"
	"path/filepath"
)

func main() {
	Serve(func(s *grpc.Server) {
		lis, err := net.Listen("tcp", "0.0.0.0:50051")
		if err != nil {
			panic(err)
		}
		if err := s.Serve(lis); err != nil {
			fmt.Printf("error serving: %v\n", err)
		}
	})
}

func Serve(servFunc func(s *grpc.Server)) {
	params := ReadSpannerParams()
	service := pb.NewUServBuilder().
		WithDefaultQueryHandlers().
		WithSpannerURI(context.Background(), params.URI()).
		WithRestOfGrpcHandlers(&RestOfImpl{params}).
		MustBuild()
	server := grpc.NewServer()

	pb.RegisterUServServer(server, service)

	servFunc(server)
}

type RestOfImpl struct {
	params SpannerParams
}

func (d *RestOfImpl) CreateTable(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	adminClient, err := admin.NewDatabaseAdminClient(ctx)
	if err != nil {
		return nil, err
	}
	op, err := adminClient.CreateDatabase(ctx, &db.CreateDatabaseRequest{
		Parent:          d.params.Parent(),
		CreateStatement: fmt.Sprintf("CREATE DATABASE %s", d.params.DatabaseId),
		ExtraStatements: []string{
			`CREATE TABLE users (
				id INT64 NOT NULL,
				name STRING(MAX) NOT NULL,
				friends BYTES(MAX) NOT NULL,
				created_on STRING(MAX) NOT NULL,
				favorite_numbers ARRAY<INT64> NOT NULL,
			) PRIMARY KEY (id)`,
		},
	})
	if err != nil {
		return nil, err
	}
	if _, err := op.Wait(ctx); err != nil {
		return nil, err
	}
	adminClient.Close()

	return &pb.Empty{}, nil
}
func (d *RestOfImpl) DropTable(ctx context.Context, req *pb.Empty) (*pb.Empty, error) {
	adminClient, err := admin.NewDatabaseAdminClient(ctx)
	if err != nil {
		return nil, err
	}
	err = adminClient.DropDatabase(ctx, &db.DropDatabaseRequest{Database: d.params.URI()})
	if err != nil {
		return nil, err
	}
	adminClient.Close()

	return &pb.Empty{}, nil
}

// using the persist lib queries to implement your own handlers.
func (d *RestOfImpl) UpdateAllNames(req *pb.Empty, stream pb.UServ_UpdateAllNamesServer) error {
	client, err := spanner.NewClient(stream.Context(), d.params.URI())
	if err != nil {
		return err
	}
	// convert our request type to a persist's type to use it in the query.
	params, err := pb.EmptyToUServPersistType(req)
	if err != nil {
		return err
	}
	// create the query using the persist type we got above (params)
	iter := client.Single().Query(stream.Context(), pl.UServGetAllUsersQuery(params))
	muts := make([]*spanner.Mutation, 0)
	err = iter.Do(func(r *spanner.Row) error {
		// get our pb.User protobuf object from a spanner row.
		user, err := pb.UserFromUServRow(r)
		if err != nil {
			return err
		}
		// turn our user object into persist's type to use in the query
		params, err := pb.UserToUServPersistType(user)
		if err != nil {
			return err
		}
		// get our mutation for this iteration's user and add it to the group
		muts = append(muts, pl.UServUpdateNameToFooQuery(params))
		return nil
	})
	if err != nil {
		return err
	}
	// apply our mutations
	if _, err := client.Apply(stream.Context(), muts); err != nil {
		return err
	}
	// get all our updated users, and stream them back to the client.
	params, _ = pb.EmptyToUServPersistType(req)
	iter = client.Single().Query(stream.Context(), pl.UServGetAllUsersQuery(params))
	err = iter.Do(func(r *spanner.Row) error {
		user, err := pb.UserFromUServRow(r)
		if err != nil {
			return err
		}
		if err := stream.Send(user); err != nil {
			return err
		}
		return nil
	})
	return err
}

type SpannerParams struct {
	ProjectId  string `json:"projectId,omitempty"`
	InstanceId string `json:"instanceId,omitempty"`
	DatabaseId string `json:"databaseId,omitempty"`
}

func (s SpannerParams) URI() string {
	return fmt.Sprintf("%s/databases/%s", s.Parent(), s.DatabaseId)
}

func (s SpannerParams) Parent() string {
	return fmt.Sprintf("projects/%s/instances/%s", s.ProjectId, s.InstanceId)
}

// need to have a struct that looks like this saved in your
// ~/.protoc-gen-persist-db.json
// {
// 		"projectId": "my-google-cloud-project-id",
// 		"instanceId": "my-google-cloud-instance-id",
// 		"databaseId": "my-google-cloud-database-id"
// }
func ReadSpannerParams() (out SpannerParams) {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	f, err := ioutil.ReadFile(filepath.Join(usr.HomeDir, "/.protoc-gen-persist-db.json"))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(f, &out); err != nil {
		panic(err)
	}
	return
}
