package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	main "github.com/tcncloud/protoc-gen-persist/test_service/user_spanner"
	"github.com/tcncloud/protoc-gen-persist/test_service/user_spanner/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"net"
	"time"
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
	main.Serve(func(s *grpc.Server) {
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
})

var _ = AfterSuite(func() {
	if testServer != nil {
		testServer.Stop()
	}
})

var _ = Describe("persist", func() {
	It("can create a table", func() {
		_, err := client.CreateTable(context.Background(), &pb.Empty{})
		Expect(err).ToNot(HaveOccurred())
	})

	It("can insert a lot of users and set their ids with before hook", func() {
		stream, err := client.InsertUsers(context.Background())
		Expect(err).To(Not(HaveOccurred()))

		for _, u := range users {
			if err := stream.Send(u); err != nil {
				Fail(err.Error())
			}
		}
		_, err = stream.CloseAndRecv()
		Expect(err).ToNot(HaveOccurred())

		retStream, err := client.GetAllUsers(context.Background(), &pb.Empty{})
		Expect(err).ToNot(HaveOccurred())

		retUsers := make([]*pb.User, 0)
		for {
			u, err := retStream.Recv()
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

	It("can select a user by id", func() {
		u, err := client.SelectUserById(context.Background(), &pb.User{Id: 0})
		Expect(err).ToNot(HaveOccurred())
		u.Id = -1
		Expect(u).To(BeEquivalentTo(users[0]))
	})

	It("can select all friends of foo", func() {
		foo, err := client.SelectUserById(context.Background(), &pb.User{Id: 0}) // foo
		Expect(err).ToNot(HaveOccurred())
		fmt.Printf("%+v\n", foo.Friends.Names)
		stream, err := client.GetFriends(context.Background(), foo.Friends)
		Expect(err).ToNot(HaveOccurred())

		friends := make([]*pb.User, 0)
		for {
			u, err := stream.Recv()
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())
			friends = append(friends, u)
		}
		Expect(friends).To(HaveLen(2))

		for _, u := range friends {
			Expect(u.Friends.Names).To(ContainElement("foo"))
		}
	})

	It("can use a client stream to update names", func() {
		stream, err := client.UpdateUserNames(context.Background())
		Expect(err).ToNot(HaveOccurred())

		for i := 0; i < len(users); i++ {
			err := stream.Send(&pb.User{Id: int64(i), Name: "zed"})
			Expect(err).ToNot(HaveOccurred())
		}
		_, err = stream.CloseAndRecv()
		Expect(err).ToNot(HaveOccurred())

		// verify changes
		retStream, err := client.GetAllUsers(context.Background(), &pb.Empty{})
		Expect(err).ToNot(HaveOccurred())

		for {
			u, err := retStream.Recv()
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(u.Name).To(Equal("zed"))
		}
	})

	It("can drop a table", func() {
		_, err := client.DropTable(context.Background(), &pb.Empty{})
		Expect(err).ToNot(HaveOccurred())
	})
})

func mustTimestamp(now time.Time) *timestamp.Timestamp {
	t, _ := ptypes.TimestampProto(now)
	return t
}
func mustNow() *timestamp.Timestamp { return mustTimestamp(time.Now()) }

var users = []*pb.User{
	&pb.User{
		Id:              -1,
		Name:            "foo",
		Friends:         &pb.Friends{Names: []string{"bar", "baz"}},
		FavoriteNumbers: []int64{1, 2, 3},
		CreatedOn:       mustNow(),
	},
	&pb.User{
		Id:              -1,
		Name:            "bar",
		Friends:         &pb.Friends{Names: []string{"foo", "baz"}},
		FavoriteNumbers: []int64{4, 5, 6},
		CreatedOn:       mustNow(),
	},
	&pb.User{
		Id:              -1,
		Name:            "baz",
		Friends:         &pb.Friends{Names: []string{"foo", "bar"}},
		FavoriteNumbers: []int64{7, 8, 9},
		CreatedOn:       mustNow(),
	},
	&pb.User{
		Id:              -1,
		Name:            "zed",
		Friends:         &pb.Friends{},
		FavoriteNumbers: []int64{1, 4, 7},
		CreatedOn:       mustNow(),
	},
}
