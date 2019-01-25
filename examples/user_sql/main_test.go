package main_test

import (
	"database/sql"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/tcncloud/protoc-gen-persist/examples/user_sql"

	"fmt"
	"io"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/tcncloud/protoc-gen-persist/examples/user_sql/pb"
	"golang.org/x/net/context"
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
		strUsers := make([]string, 0)
		for _, u := range users {
			strUsers = append(strUsers, proto.MarshalTextString(u))
		}

		for _, u := range retUsers {
			Expect(strUsers).To(ContainElement(proto.MarshalTextString(u)))
		}
	})

	It("can select a user by id", func() {
		u, err := client.SelectUserById(context.Background(), &pb.User{Id: 0})
		Expect(err).ToNot(HaveOccurred())
		Expect(u.Id).To(BeEquivalentTo(0))
		u.Id = -1
		Expect(proto.MarshalTextString(u)).To(BeEquivalentTo(proto.MarshalTextString(users[0])))
	})

	It("can select all friends of foo", func() {
		foo, err := client.SelectUserById(context.Background(), &pb.User{Id: 0}) // foo
		Expect(err).ToNot(HaveOccurred())
		stream, err := client.GetFriends(context.Background(), &pb.FriendsReq{
			Names: &pb.SliceStringParam{Slice: foo.Friends.Names},
		})
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

	It("can use a bidirectional stream to update names", func() {
		stream, err := client.UpdateUserNames(context.Background())
		Expect(err).ToNot(HaveOccurred())

		for i := 0; i < len(users); i++ {
			err := stream.Send(&pb.User{Id: int64(i), Name: "zed"})
			Expect(err).ToNot(HaveOccurred())
		}
		err = stream.CloseSend()
		Expect(err).ToNot(HaveOccurred())
		for {
			u, err := stream.Recv()
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(u.Name).To(Equal("zed"))
		}
	})

	It("can change all names to foo", func() {
		stream, err := client.UpdateAllNames(context.Background(), &pb.Empty{})
		Expect(err).ToNot(HaveOccurred())
		var resps int
		for {
			u, err := stream.Recv()
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())
			Expect(u.Name).To(Equal("foo"))
			resps++
		}
		Expect(resps).To(BeNumerically(">", 0))
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
		Id:        -1,
		Name:      "foo",
		Friends:   &pb.Friends{Names: []string{"bar", "baz"}},
		CreatedOn: mustNow(),
		Id2: 35,
	},
	&pb.User{
		Id:        -1,
		Name:      "bar",
		Friends:   &pb.Friends{Names: []string{"foo", "baz"}},
		CreatedOn: mustNow(),
		Id2: 35,
	},
	&pb.User{
		Id:        -1,
		Name:      "baz",
		Friends:   &pb.Friends{Names: []string{"foo", "bar"}},
		CreatedOn: mustNow(),
		Id2: 35,
	},
	&pb.User{
		Id:        -1,
		Name:      "zed",
		Friends:   &pb.Friends{},
		CreatedOn: mustNow(),
		Id2: 35,
	},
}

func Serve(servFunc func(s *grpc.Server)) {
	conn, err := sql.Open("postgres", "user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	opts := pb.OptsUServ(&HooksImpl{}, &MappingImpl{})
	handlers := &RestOfImpl{
		DB:      conn,
		QUERIES: pb.QueriesUServ(opts),
	}
	service := pb.ImplUServ(conn, handlers, opts)

	server := grpc.NewServer()

	pb.RegisterUServServer(server, service)

	servFunc(server)
}
