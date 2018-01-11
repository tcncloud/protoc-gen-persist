// This file is generated by protoc-gen-persist
// Source File: pb/user.proto
// DO NOT EDIT !
package pb

import (
	io "io"

	spanner "cloud.google.com/go/spanner"
	proto "github.com/golang/protobuf/proto"
	persist_lib "github.com/tcncloud/protoc-gen-persist/test_service/user_spanner/pb/persist_lib"
	context "golang.org/x/net/context"
	codes "google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

type UServImpl struct {
	PERSIST   *persist_lib.UServMethodReceiver
	FORWARDED RestOfUServHandlers
}
type RestOfUServHandlers interface {
	CreateTable(ctx context.Context, req *Empty) (*Empty, error)
	DropTable(ctx context.Context, req *Empty) (*Empty, error)
}
type UServImplBuilder struct {
	err           error
	rest          RestOfUServHandlers
	queryHandlers *persist_lib.UServQueryHandlers
	i             *UServImpl
	db            spanner.Client
}

func NewUServBuilder() *UServImplBuilder {
	return &UServImplBuilder{i: &UServImpl{}}
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
	accessor := persist_lib.NewSpannerClientGetter(&b.db)
	queryHandlers := &persist_lib.UServQueryHandlers{
		InsertUsersHandler:     persist_lib.DefaultInsertUsersHandler(accessor),
		GetAllUsersHandler:     persist_lib.DefaultGetAllUsersHandler(accessor),
		SelectUserByIdHandler:  persist_lib.DefaultSelectUserByIdHandler(accessor),
		UpdateUserNamesHandler: persist_lib.DefaultUpdateUserNamesHandler(accessor),
		GetFriendsHandler:      persist_lib.DefaultGetFriendsHandler(accessor),
	}
	b.queryHandlers = queryHandlers
	return b
}

// set the custom handlers you want to use in the handlers
// this method will make sure to use a default handler if
// the handler is nil.
func (b *UServImplBuilder) WithNilAsDefaultQueryHandlers(p *persist_lib.UServQueryHandlers) *UServImplBuilder {
	accessor := persist_lib.NewSpannerClientGetter(&b.db)
	if p.InsertUsersHandler == nil {
		p.InsertUsersHandler = persist_lib.DefaultInsertUsersHandler(accessor)
	}
	if p.GetAllUsersHandler == nil {
		p.GetAllUsersHandler = persist_lib.DefaultGetAllUsersHandler(accessor)
	}
	if p.SelectUserByIdHandler == nil {
		p.SelectUserByIdHandler = persist_lib.DefaultSelectUserByIdHandler(accessor)
	}
	if p.UpdateUserNamesHandler == nil {
		p.UpdateUserNamesHandler = persist_lib.DefaultUpdateUserNamesHandler(accessor)
	}
	if p.GetFriendsHandler == nil {
		p.GetFriendsHandler = persist_lib.DefaultGetFriendsHandler(accessor)
	}
	b.queryHandlers = p
	return b
}
func (b *UServImplBuilder) WithSpannerClient(c *spanner.Client) *UServImplBuilder {
	b.db = *c
	return b
}
func (b *UServImplBuilder) WithSpannerURI(ctx context.Context, uri string) *UServImplBuilder {
	cli, err := spanner.NewClient(ctx, uri)
	b.err = err
	b.db = *cli
	return b
}
func (b *UServImplBuilder) Build() (*UServImpl, error) {
	if b.err != nil {
		return nil, b.err
	}
	b.i.PERSIST = &persist_lib.UServMethodReceiver{Handlers: *b.queryHandlers}
	b.i.FORWARDED = b.rest
	return b.i, nil
}
func (b *UServImplBuilder) MustBuild() *UServImpl {
	s, err := b.Build()
	if err != nil {
		panic("error in builder: " + err.Error())
	}
	return s
}

func (s *UServImpl) CreateTable(ctx context.Context, req *Empty) (*Empty, error) {
	return s.FORWARDED.CreateTable(ctx, req)
}

func (s *UServImpl) InsertUsers(stream UServ_InsertUsersServer) error {
	var err error
	_ = err
	res := Empty{}
	feed, stop := s.PERSIST.InsertUsers(stream.Context())
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return gstatus.Errorf(codes.Unknown, "error receiving request: %v", err)
		}
		beforeRes, err := IncId(req)
		if err != nil {
			return gstatus.Errorf(codes.Unknown, "error in before hook: %v", err)
		} else if beforeRes != nil {
			continue
		}
		params := &persist_lib.UserForUServ{}
		err = func() error {
			// set 'User.id' in params
			params.Id = req.Id
			// set 'User.name' in params
			params.Name = req.Name
			// set 'User.friends' in params
			if req.Friends == nil {
				req.Friends = new(Friends)
			}
			{
				raw, err := proto.Marshal(req.Friends)
				if err != nil {
					return err
				}
				params.Friends = raw
			}
			// set 'User.created_on' in params
			if params.CreatedOn, err = (TimeString{}).ToSpanner(req.CreatedOn).SpannerValue(); err != nil {
				return err
			}
			// set 'User.favorite_numbers' in params
			params.FavoriteNumbers = req.FavoriteNumbers
			return nil
		}()
		if err != nil {
			return err
		}
		feed(params)
	}
	row, err := stop()
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error receiving result row: %v", err)
	}
	if row != nil {
		err = func() error {
			return nil
		}()
	}
	if err := stream.SendAndClose(&res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
	}
	return nil
}

func (s *UServImpl) GetAllUsers(req *Empty, stream UServ_GetAllUsersServer) error {
	var err error
	_ = err
	params := &persist_lib.EmptyForUServ{}
	err = func() error {
		return nil
	}()
	if err != nil {
		return err
	}
	var iterErr error
	err = s.PERSIST.GetAllUsers(stream.Context(), params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res := User{}
		err = func() error {
			var Id_ int64
			{
				local := &spanner.NullInt64{}
				if err := row.ColumnByName("id", local); err != nil {
					return err
				}
				if local.Valid {
					Id_ = local.Int64
				}
				res.Id = Id_
			}
			var Name_ string
			{
				local := &spanner.NullString{}
				if err := row.ColumnByName("name", local); err != nil {
					return err
				}
				if local.Valid {
					Name_ = local.StringVal
				}
				res.Name = Name_
			}
			var Friends_ []byte
			if err := row.ColumnByName("friends", &Friends_); err != nil {
				return err
			}
			{
				local := new(Friends)
				if err := proto.Unmarshal(Friends_, local); err != nil {
					return err
				}
				res.Friends = local
			}
			var CreatedOn_ = new(spanner.GenericColumnValue)
			if err := row.ColumnByName("created_on", CreatedOn_); err != nil {
				return err
			}
			{
				local := &TimeString{}
				if err := local.SpannerScan(CreatedOn_); err != nil {
					return err
				}
				res.CreatedOn = local.ToProto()
			}
			var FavoriteNumbers_ []int64
			{
				local := make([]spanner.NullInt64, 0)
				if err := row.ColumnByName("favorite_numbers", &local); err != nil {
					return err
				}
				for _, l := range local {
					if l.Valid {
						FavoriteNumbers_ = append(FavoriteNumbers_, l.Int64)
						res.FavoriteNumbers = FavoriteNumbers_
					}
				}
			}
			return nil
		}()
		if err != nil {
			iterErr = err
			return
		}
		if err := stream.Send(&res); err != nil {
			iterErr = gstatus.Errorf(codes.Unknown, "error during iteration: %v", err)
		}
	})
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error during iteration: %v", err)
	} else if iterErr != nil {
		return iterErr
	}
	return nil
}

func (s *UServImpl) SelectUserById(ctx context.Context, req *User) (*User, error) {
	var err error
	var res = User{}
	_ = err
	_ = res
	params := &persist_lib.UserForUServ{}
	err = func() error {
		// set 'User.id' in params
		params.Id = req.Id
		// set 'User.name' in params
		params.Name = req.Name
		// set 'User.friends' in params
		if req.Friends == nil {
			req.Friends = new(Friends)
		}
		{
			raw, err := proto.Marshal(req.Friends)
			if err != nil {
				return err
			}
			params.Friends = raw
		}
		// set 'User.created_on' in params
		if params.CreatedOn, err = (TimeString{}).ToSpanner(req.CreatedOn).SpannerValue(); err != nil {
			return err
		}
		// set 'User.favorite_numbers' in params
		params.FavoriteNumbers = req.FavoriteNumbers
		return nil
	}()
	if err != nil {
		return nil, err
	}
	var iterErr error
	err = s.PERSIST.SelectUserById(ctx, params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res = User{}
		err = func() error {
			var Id_ int64
			{
				local := &spanner.NullInt64{}
				if err := row.ColumnByName("id", local); err != nil {
					return err
				}
				if local.Valid {
					Id_ = local.Int64
				}
				res.Id = Id_
			}
			var Name_ string
			{
				local := &spanner.NullString{}
				if err := row.ColumnByName("name", local); err != nil {
					return err
				}
				if local.Valid {
					Name_ = local.StringVal
				}
				res.Name = Name_
			}
			var Friends_ []byte
			if err := row.ColumnByName("friends", &Friends_); err != nil {
				return err
			}
			{
				local := new(Friends)
				if err := proto.Unmarshal(Friends_, local); err != nil {
					return err
				}
				res.Friends = local
			}
			var CreatedOn_ = new(spanner.GenericColumnValue)
			if err := row.ColumnByName("created_on", CreatedOn_); err != nil {
				return err
			}
			{
				local := &TimeString{}
				if err := local.SpannerScan(CreatedOn_); err != nil {
					return err
				}
				res.CreatedOn = local.ToProto()
			}
			var FavoriteNumbers_ []int64
			{
				local := make([]spanner.NullInt64, 0)
				if err := row.ColumnByName("favorite_numbers", &local); err != nil {
					return err
				}
				for _, l := range local {
					if l.Valid {
						FavoriteNumbers_ = append(FavoriteNumbers_, l.Int64)
						res.FavoriteNumbers = FavoriteNumbers_
					}
				}
			}
			return nil
		}()
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
	return &res, nil
}

func (s *UServImpl) UpdateUserNames(stream UServ_UpdateUserNamesServer) error {
	var err error
	_ = err
	res := Empty{}
	feed, stop := s.PERSIST.UpdateUserNames(stream.Context())
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return gstatus.Errorf(codes.Unknown, "error receiving request: %v", err)
		}
		params := &persist_lib.UserForUServ{}
		err = func() error {
			// set 'User.id' in params
			params.Id = req.Id
			// set 'User.name' in params
			params.Name = req.Name
			// set 'User.friends' in params
			if req.Friends == nil {
				req.Friends = new(Friends)
			}
			{
				raw, err := proto.Marshal(req.Friends)
				if err != nil {
					return err
				}
				params.Friends = raw
			}
			// set 'User.created_on' in params
			if params.CreatedOn, err = (TimeString{}).ToSpanner(req.CreatedOn).SpannerValue(); err != nil {
				return err
			}
			// set 'User.favorite_numbers' in params
			params.FavoriteNumbers = req.FavoriteNumbers
			return nil
		}()
		if err != nil {
			return err
		}
		feed(params)
	}
	row, err := stop()
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error receiving result row: %v", err)
	}
	if row != nil {
		err = func() error {
			return nil
		}()
	}
	if err := stream.SendAndClose(&res); err != nil {
		return gstatus.Errorf(codes.Unknown, "error sending back response: %v", err)
	}
	return nil
}

func (s *UServImpl) GetFriends(req *Friends, stream UServ_GetFriendsServer) error {
	var err error
	_ = err
	params := &persist_lib.FriendsForUServ{}
	err = func() error {
		// set 'Friends.names' in params
		params.Names = req.Names
		return nil
	}()
	if err != nil {
		return err
	}
	var iterErr error
	err = s.PERSIST.GetFriends(stream.Context(), params, func(row *spanner.Row) {
		if row == nil { // there was no return data
			return
		}
		res := User{}
		err = func() error {
			var Id_ int64
			{
				local := &spanner.NullInt64{}
				if err := row.ColumnByName("id", local); err != nil {
					return err
				}
				if local.Valid {
					Id_ = local.Int64
				}
				res.Id = Id_
			}
			var Name_ string
			{
				local := &spanner.NullString{}
				if err := row.ColumnByName("name", local); err != nil {
					return err
				}
				if local.Valid {
					Name_ = local.StringVal
				}
				res.Name = Name_
			}
			var Friends_ []byte
			if err := row.ColumnByName("friends", &Friends_); err != nil {
				return err
			}
			{
				local := new(Friends)
				if err := proto.Unmarshal(Friends_, local); err != nil {
					return err
				}
				res.Friends = local
			}
			var CreatedOn_ = new(spanner.GenericColumnValue)
			if err := row.ColumnByName("created_on", CreatedOn_); err != nil {
				return err
			}
			{
				local := &TimeString{}
				if err := local.SpannerScan(CreatedOn_); err != nil {
					return err
				}
				res.CreatedOn = local.ToProto()
			}
			var FavoriteNumbers_ []int64
			{
				local := make([]spanner.NullInt64, 0)
				if err := row.ColumnByName("favorite_numbers", &local); err != nil {
					return err
				}
				for _, l := range local {
					if l.Valid {
						FavoriteNumbers_ = append(FavoriteNumbers_, l.Int64)
						res.FavoriteNumbers = FavoriteNumbers_
					}
				}
			}
			return nil
		}()
		if err != nil {
			iterErr = err
			return
		}
		if err := stream.Send(&res); err != nil {
			iterErr = gstatus.Errorf(codes.Unknown, "error during iteration: %v", err)
		}
	})
	if err != nil {
		return gstatus.Errorf(codes.Unknown, "error during iteration: %v", err)
	} else if iterErr != nil {
		return iterErr
	}
	return nil
}

func (s *UServImpl) DropTable(ctx context.Context, req *Empty) (*Empty, error) {
	return s.FORWARDED.DropTable(ctx, req)
}