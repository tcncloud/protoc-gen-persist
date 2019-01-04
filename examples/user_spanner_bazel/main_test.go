package main_test

import (
	"testing"
  "io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"net"
	"time"

  spanner "cloud.google.com/go/spanner"
	admin "cloud.google.com/go/spanner/admin/database/apiv1"
	ptypess "github.com/golang/protobuf/ptypes"
	timeystamp "github.com/golang/protobuf/ptypes/timestamp"
	main "github.com/tcncloud/protoc-gen-persist/examples/user_spanner_bazel"
	pb "github.com/tcncloud/protoc-gen-persist/examples/user_spanner_bazel/pb"
	"golang.org/x/net/context"
	db "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/grpc"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var (
	testServer *grpc.Server
	client     pb.UServClient
)

var _ = BeforeSuite(func() {
	Serve(func(s *grpc.Server) {
		lis, err := net.Listen("tcp", "0.0.0.0:50051")
		if err != nil {
			Fail("could not register listener: " + err.Error())
			return
		}
		testServer = s
		go func() {
			if err := s.Serve(lis); err != nil {
				fmt.Printf("failed? %v", err)
			}
		}()
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			Fail("could not create the client: " + err.Error())
		}
		client = pb.NewUServClient(conn)
	})
	err := CreateTable(context.Background(), main.ReadSpannerParams())
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	if testServer != nil {
		testServer.Stop()

		err := DropTable(context.Background(), main.ReadSpannerParams())
		Expect(err).ToNot(HaveOccurred())
	}
})

var _ = Describe("persist", func() {
	It("can create a table", func() {

	})

	It("can insert a lot of users and set their ids with before hook", func() {
		stream, err := client.InsertUsers(context.Background())
		Expect(err).To(Not(HaveOccurred()))

    fmt.Println("inserted users")
		for _, u := range users {
			if err := stream.Send(u); err != nil {
				Fail(err.Error())
			}
		}
		_, err = stream.CloseAndRecv()
		Expect(err).ToNot(HaveOccurred())

		retStream, err := client.GetAllUsers(context.Background(), &pb.Empty{})
		Expect(err).ToNot(HaveOccurred())
    fmt.Println("retStream", retStream)
    fmt.Println("err", err)

		retUsers := make([]*pb.User, 0)
		for {
      // TODO start here. why is this getting an io.eof
			u, err := retStream.Recv()
      fmt.Println("u: ", u)
      fmt.Println("err: ", err)
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(u.Id).ToNot(Equal(-1))
			u.Id = -1
			retUsers = append(retUsers, u)
		}
		Expect(retUsers).To(HaveLen(len(users)))
		for _, u := range retUsers {
			Expect(users).To(ContainElement(BeEquivalentTo(u)))
		}
	})

	// PIt("can select a user by id", func() {
	// 	u, err := client.SelectUserById(context.Background(), &pb.User{Id: 0})
	// 	Expect(err).ToNot(HaveOccurred())
	// 	u.Id = -1
	// 	Expect(u).To(BeEquivalentTo(users[0]))
	// })

	// PIt("can select all friends of foo", func() {
	// 	foo, err := client.SelectUserById(context.Background(), &pb.User{Id: 0}) // foo
	// 	Expect(err).ToNot(HaveOccurred())
	// 	stream, err := client.GetFriends(context.Background(), foo.Friends)
	// 	Expect(err).ToNot(HaveOccurred())

	// 	friends := make([]*pb.User, 0)
	// 	for {
	// 		u, err := stream.Recv()
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		Expect(err).ToNot(HaveOccurred())
	// 		friends = append(friends, u)
	// 	}
	// 	Expect(friends).To(HaveLen(2))

	// 	for _, u := range friends {
	// 		Expect(u.Friends.Names).To(ContainElement("foo"))
	// 	}
	// })

	// PIt("can use a client stream to update names", func() {
	// 	stream, err := client.UpdateUserNames(context.Background())
	// 	Expect(err).ToNot(HaveOccurred())

	// 	for i := 0; i < len(users); i++ {
	// 		err := stream.Send(&pb.User{Id: int64(i), Name: "zed"})
	// 		Expect(err).ToNot(HaveOccurred())
	// 	}
	// 	_, err = stream.CloseAndRecv()
	// 	Expect(err).ToNot(HaveOccurred())

	// 	// verify changes
	// 	retStream, err := client.GetAllUsers(context.Background(), &pb.Empty{})
	// 	Expect(err).ToNot(HaveOccurred())

	// 	for {
	// 		u, err := retStream.Recv()
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		Expect(err).ToNot(HaveOccurred())
	// 		Expect(u.Name).To(Equal("zed"))
	// 	}
	// })
})

func mustTimestamp(now time.Time) *timeystamp.Timestamp {
	t, _ := ptypess.TimestampProto(now)
	return t
}
func mustNow() *timeystamp.Timestamp { return mustTimestamp(time.Now()) }

var users = []*pb.User{
	&pb.User{
		Id:              -1,
		Name:            "foo",
		Friends:         &pb.Friends{Names: []string{"bar", "baz"}},
		CreatedOn:       mustNow(),
	},
	&pb.User{
		Id:              -1,
		Name:            "bar",
		Friends:         &pb.Friends{Names: []string{"foo", "baz"}},
		CreatedOn:       mustNow(),
	},
	&pb.User{
		Id:              -1,
		Name:            "baz",
		Friends:         &pb.Friends{Names: []string{"foo", "bar"}},
		CreatedOn:       mustNow(),
	},
	&pb.User{
		Id:              -1,
		Name:            "zed",
		Friends:         &pb.Friends{},
		CreatedOn:       mustNow(),
	},
}

func Serve(servFunc func(s *grpc.Server)) {
  params := main.ReadSpannerParams()
  ctx := context.Background()
  conn, err := spanner.NewClient(ctx, params.URI())
  if err != nil {
    fmt.Printf("error connecting to db: %v\n", err)
    return
  }
  // defer conn.Close()

  service := pb.UServPersistImpl(conn, pb.UServ_ImplOpts{
    HOOKS: &main.HooksImpl{},
    MAPPINGS: &main.MappingImpl{},
    HANDLERS: &main.RestOfImpl{},
  })

  server := grpc.NewServer()

  pb.RegisterUServServer(server, service)

	servFunc(server)
}

func CreateTable(ctx context.Context, params main.SpannerParams) error {
	adminClient, err := admin.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	op, err := adminClient.CreateDatabase(ctx, &db.CreateDatabaseRequest{
		Parent:          params.Parent(),
		CreateStatement: fmt.Sprintf("CREATE DATABASE %s", params.DatabaseId),
    // TODO this probably needs to be an array of bytes because it is MULTIPLE friends, and not just one
		ExtraStatements: []string{
			`CREATE TABLE users (
				id INT64 NOT NULL,
				name STRING(MAX) NOT NULL,
				friends BYTES(MAX) NOT NULL,
				created_on STRING(MAX) NOT NULL,
			) PRIMARY KEY (id)`,
		},
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	// adminClient.Close()

	return nil
}
func DropTable(ctx context.Context, params main.SpannerParams) error {
	adminClient, err := admin.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	err = adminClient.DropDatabase(ctx, &db.DropDatabaseRequest{Database: params.URI()})
	if err != nil {
		return err
	}
	// adminClient.Close()

	return nil
}
