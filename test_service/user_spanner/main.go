package main

import (
	admin "cloud.google.com/go/spanner/admin/database/apiv1"
	"encoding/json"
	"fmt"
	"github.com/tcncloud/protoc-gen-persist/test_service/user_spanner/pb"
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
