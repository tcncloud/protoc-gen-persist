load("@io_bazel_rules_go//go:def.bzl", "go_library")

def get_outputs(ctx, proto_files, file_names):
    outputs = []

    for f in proto_files:
        base = f.basename[:-len(".proto")]

        for file_name in file_names:
            new_file_name = file_name.format(basename = base)
            new_file = ctx.actions.declare_file(new_file_name)
            outputs += [new_file]

    return outputs

# TODO: Make sure this works for internal dependencies
def get_proto_directories(proto_deps):
    proto_path = "--proto_path=. \\"
    for dep in proto_deps:
        for f in dep.files:
            path_pieces = f.dirname.split("/")
            proto_path += "--proto_path=%s/%s" % (path_pieces[0], path_pieces[1])

    return proto_path

def get_proto_file_paths(proto_files):
    proto_paths = ""
    for f in proto_files:
        proto_paths += "%s \\" % f.path
    
    return proto_paths

def _persist_lib_pkg_compile_impl(ctx):
    file_names = [
        "persist_lib/{basename}_query_handlers.persist.go",
        "persist_lib/{basename}_queries.persist.go",
    ]
    outputs = get_outputs(ctx, ctx.attr.proto.files, file_names)
    outputs += [ctx.actions.declare_file("persist_lib/pkg_level_definitions.persist.go")]

    proto_paths = get_proto_file_paths(ctx.attr.proto.files)
    proto_directories = get_proto_directories(ctx.attr.proto_deps)

    ctx.actions.run_shell(
        inputs = ctx.files.proto + ctx.files.proto_deps + [ctx.file._persist_plugin],
        outputs = outputs,
        progress_message = "Test message",
        command = "protoc \
        --plugin=protoc-gen-persist=%s \
        --persist_out=persist_root=%s:%s \
        %s \
        %s" % (
            ctx.file._persist_plugin.path,
            ctx.attr.persist_root,
            ctx.var["GENDIR"],
            proto_directories,
            proto_paths,
        ),
    )

    return DefaultInfo(
        files=depset(outputs),
    )

persist_lib_pkg_compile = rule(
    attrs = {
        "persist_root": attr.string(),
        "proto": attr.label(
            mandatory = True,
            allow_files = True,
        ),
        "proto_deps": attr.label_list(),
        "_persist_plugin": attr.label(
            default = Label("//:protoc-gen-persist"),
            executable = True,
            allow_single_file = True,
            cfg = "host",  #TODO: figure out what this means
        ),
    },
    output_to_genfiles = True,
    implementation = _persist_lib_pkg_compile_impl,
)

def _lib_pkg_compile_impl(ctx):
    file_names = [
        "{basename}.persist.go",
        "{basename}.pb.go",
    ]
    outputs = get_outputs(ctx, ctx.attr.proto.files, file_names)
    proto_paths = get_proto_file_paths(ctx.attr.proto.files)
    proto_directories = get_proto_directories(ctx.attr.proto_deps)

    ctx.actions.run_shell(
        inputs = ctx.files.proto + ctx.files.proto_deps + [ctx.file._go_plugin] + [ctx.file._persist_plugin],
        outputs = outputs,
        progress_message = "Test message",
        command = "protoc \
        --plugin=protoc-gen-persist=%s \
        --persist_out=persist_root=%s:%s \
        --plugin=protoc-gen-go=%s \
        --go_out=plugins=grpc:%s \
        %s \
        %s" % (
            ctx.file._persist_plugin.path,
            ctx.attr.persist_root,
            ctx.var["GENDIR"],
            ctx.file._go_plugin.path,
            ctx.var["GENDIR"],
            proto_directories,
            proto_paths,
        ),
    )

    return DefaultInfo(
        files=depset(outputs),
    )

lib_pkg_compile = rule(
    attrs = {
        "persist_root": attr.string(),
        "proto": attr.label(
            mandatory = True,
            allow_files = True,
        ),
        "proto_deps": attr.label_list(),
        "_persist_plugin": attr.label(
            default = Label("//:protoc-gen-persist"),
            executable = True,
            allow_single_file = True,
            cfg = "host",  #TODO: figure out what this means
        ),
        "_go_plugin": attr.label(
            default = Label("@com_github_golang_protobuf//protoc-gen-go:protoc-gen-go"),
            executable = True,
            allow_single_file = True,
            cfg = "host",  #TODO: figure out what this means
        ),
    },
    output_to_genfiles = True,
    implementation = _lib_pkg_compile_impl,
)

def persist_go_library(
    name,
    go_srcs = [],
    go_lib_deps = [],
    importpath = None,
    proto = [],
    proto_deps = [],
    visibility = None,
    **kwargs):

    # Generate files in persist lib package
    persist_lib = name + "_persist_lib"
    persist_lib_pkg_compile(
        name = persist_lib,
        persist_root = importpath,
        proto = proto,
        proto_deps = proto_deps,
    )

    # Generate files in lib package
    lib = name + "_lib"
    lib_pkg_compile(
        name = lib,
        persist_root = importpath,
        proto = proto,
        proto_deps = proto_deps,
    )
    
    # Compile go library from persist lib files
    go_lib_1 = name + "_go_lib_1"
    go_library(
        name = go_lib_1,
        srcs = [":" + persist_lib],
        visibility = visibility,
        deps = go_lib_deps,
        importpath = importpath + "/persist_lib",
    )

    # Compile go library from lib files
    go_library(
        name = name,
        srcs = [":" + lib] + go_srcs,
        visibility = visibility,
        deps = depset(go_lib_deps + [":" + go_lib_1]),
        importpath = importpath,
    )