go_mod_download(
    name = "protobuf_dl",
    module = "google.golang.org/protobuf",
    version = "v1.28.1",  # check https://github.com/grpc/grpc-go/tree/master/cmd/protoc-gen-go-grpc
    # in go.mod what's the protobuf version supported by the grpc version you're using
)

go_module(
    name = "protobuf",
    download = ":protobuf_dl",
    install = [
        "compiler/protogen",
        "encoding/protojson",
        "encoding/prototext",
        "encoding/protowire",
        "internal/descfmt",
        "internal/descopts",
        "internal/detrand",
        "internal/encoding/defval",
        "internal/encoding/json",
        "internal/encoding/messageset",
        "internal/encoding/tag",
        "internal/encoding/text",
        "internal/errors",
        "internal/filedesc",
        "internal/filetype",
        "internal/flags",
        "internal/genid",
        "internal/impl",
        "internal/order",
        "internal/pragma",
        "internal/set",
        "internal/strs",
        "internal/version",
        "proto",
        "reflect/protodesc",
        "reflect/protoreflect",
        "reflect/protoregistry",
        "runtime/protoiface",
        "runtime/protoimpl",
        "types/descriptorpb",
        "types/dynamicpb",
        "types/known/anypb",
        "types/known/durationpb",
        "types/known/emptypb",
        "types/known/fieldmaskpb",
        "types/known/structpb",
        "types/known/timestamppb",
        "types/known/wrapperspb",
        "types/pluginpb",
    ],
    module = "google.golang.org/protobuf",
    visibility = ["PUBLIC"],
)

go_module(
    name = "protoc-gen-go",
    binary = True,
    download = ":protobuf_dl",
    install = ["cmd/protoc-gen-go"],
    module = "google.golang.org/protobuf",
    visibility = ["PUBLIC"],
    deps = [":protobuf"],
)

go_module(
    name = "logrus",
    module = "github.com/sirupsen/logrus",
    version = "v1.9.0",
    visibility = ["PUBLIC"],
    deps = [":sys"],
)

go_module(
    name = "sys",
    install = [
        "execabs",
        "internal/unsafeheader",
        "unix",
    ],
    module = "golang.org/x/sys",
    version = "v0.7.0",
    visibility = ["PUBLIC"],
)

go_module(
    name = "tools",
    install = [
        "go/ast/astutil",
        "imports",
        "internal/event",
        "internal/event/core",
        "internal/event/keys",
        "internal/event/label",
        "internal/gocommand",
        "internal/gopathwalk",
        "internal/imports",
        "internal/typeparams",
        "internal/fastwalk",
    ],
    module = "golang.org/x/tools",
    version = "v0.8.0",
    visibility = ["PUBLIC"],
    deps = [
        ":mod",
        ":sys",
    ],
)

go_module(
    name = "mod",
    install = [
        "internal/lazyregexp",
        "module",
        "semver",
    ],
    module = "golang.org/x/mod",
    version = "v0.10.0",
    visibility = ["PUBLIC"],
)
