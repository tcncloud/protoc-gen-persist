load(
    "@io_bazel_rules_go//go:def.bzl",
    "GoLibrary",
    "go_library",
)
load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)
load(
    "@io_bazel_rules_go//proto:compiler.bzl",
    "proto_path",
)

def get_outputs(ctx, dynamic_file_names, static_file_names):
    outputs = []
    for src in ctx.attr.proto.proto.direct_sources:
        base = src.basename[:-len(".proto")]

        if (base.endswith("_p") == True):
            for file_name in dynamic_file_names:
                new_file_name = "%s~/%s/%s" % (ctx.attr.name, ctx.attr.importpath, file_name)
                new_file_name = new_file_name.format(basename = base)
                outputs += [ctx.actions.declare_file(new_file_name)]

    if (len(outputs) > 0):
        for file_name in static_file_names:
            new_file_name = "%s~/%s/%s" % (ctx.attr.name, ctx.attr.importpath, file_name)
            outputs += [ctx.actions.declare_file(new_file_name)]

    return outputs    

def get_proto_file_paths(proto_files):
    proto_paths = []
    for f in proto_files:
        proto_paths.append(proto_path(f))
    
    return proto_paths

def combine_inputs(*inputs):
    combined = depset()
    for i in inputs:
        combined += i
    return combined

def prepare_args(ctx, outputs):
    proto = ctx.attr.proto.proto
    outpath = outputs[0].dirname[:-len(ctx.attr.importpath)]

    if (ctx.attr.importpath.endswith("/persist_lib") == True):
        base_importpath = ctx.attr.importpath[:-len("/persist_lib")]
    else:
        base_importpath = ctx.attr.importpath

    args = ctx.actions.args()

    args.add([
        "--protoc", ctx.file._protoc,
        "--out_path", outpath,
        "--plugin", ctx.file._persist_plugin,
    ])

    # Add descriptor set
    args.add(proto.transitive_descriptor_sets, before_each = "--descriptor_set")
    
    # Add outputs
    args.add(outputs, before_each = "--expected")

    # Add imports
    imports = []
    imports.append(["{}={}".format(proto_path(src), base_importpath) for src in ctx.attr.proto.proto.direct_sources])
    args.add(combine_inputs(*imports), before_each = "--import")

    # Add proto paths (the path to the proto that was seen when the descriptor file was generated)
    proto_paths = get_proto_file_paths(proto.direct_sources)
    args.add(proto_paths)

    return args

def _persist_proto_compile_impl(ctx):
    if (ctx.attr.type == "persist_lib"):
        if (ctx.attr.importpath.endswith("persist_lib") == False):
            fail("The importpath for persist_lib_go_library must end with persist_lib")

        dynamic_file_names = [
            "{basename}_query_handlers.persist.go",
            "{basename}_queries.persist.go",
        ]
        static_file_names = [
            "pkg_level_definitions.persist.go",
        ]
    else:
        dynamic_file_names = ["{basename}.persist.go"]
        static_file_names = []

    outputs = get_outputs(ctx, dynamic_file_names, static_file_names)

    if (len(outputs) == 0):
        fail("Only filenames suffixed with _p will generate persist.go files. Please check your proto filenames.")

    args = prepare_args(ctx, outputs)

    ctx.actions.run(
        inputs = combine_inputs(
            [ctx.file._persist_plugin],
            [ctx.file._protoc],
            [ctx.file._go_protoc],
            ctx.attr.proto.proto.transitive_descriptor_sets
        ),
        outputs = outputs,
        progress_message = "Generating into %s" % outputs[0].dirname,
        mnemonic = "PersistProtocGen",
        executable = ctx.file._go_protoc,
        arguments = [args],
    )

    return DefaultInfo(
        files=depset(outputs),
    )

persist_proto_compile = rule(
    attrs = {
        "importpath": attr.string(
            mandatory = True,
        ),
        "type": attr.string(
            mandatory = True,
        ),
        "proto": attr.label(
            mandatory = True,
            allow_files = True,
            single_file = True,
            providers = ["proto"],
        ),
        "_persist_plugin": attr.label(
            default = Label("//:protoc-gen-persist"),
            executable = True,
            allow_files = True,
            single_file = True,
            cfg = "host",
        ),
        "_protoc": attr.label(
            allow_files = True,
            single_file = True,
            executable = True,
            cfg = "host",
            default = Label("@com_google_protobuf//:protoc"),
        ),
        "_go_protoc": attr.label(
            allow_files = True,
            single_file = True,
            executable = True,
            cfg = "host",
            default = Label("@io_bazel_rules_go//go/tools/builders:go-protoc"),
        ),
    },
    # output_to_genfiles = True,
    implementation = _persist_proto_compile_impl,
)

def persist_lib_go_library(
    name,
    go_lib_deps = [],
    importpath = None,
    visibility = None,
    proto = None,
    **kwargs):

    # Generate files in persist lib package
    persist_lib = name + "_files"
    persist_proto_compile(
        name = persist_lib,
        type = "persist_lib",
        importpath = importpath,
        proto = proto,
    )

    # Compile go library from persist lib files
    go_library(
        name = name,
        srcs = [":" + persist_lib],
        visibility = visibility,
        deps = go_lib_deps,
        importpath = importpath,
    )

def persist_go_library(
    name,
    go_srcs = [],
    go_lib_deps = [],
    importpath = None,
    visibility = None,
    proto = None,
    **kwargs):

    # Generate persist.go files
    persist_files = name + "_files"
    persist_proto_compile(
        name = persist_files,
        type = "persist_files",
        importpath = importpath,
        proto = proto,
    )

    # Generate pb.go files
    go_lib = name + "_go_lib"
    go_proto_library(
        name = go_lib,
        compilers = ["@io_bazel_rules_go//proto:go_grpc"],
        importpath = importpath,
        proto = proto,
        visibility = ["//visibility:public"],
        deps = go_lib_deps,
    )

    # Compile go library from lib files
    go_library(
        name = name,
        srcs = [":" + persist_files] + go_srcs,
        embed = [":" + go_lib],
        importpath = importpath,
        visibility = visibility,
        deps = go_lib_deps,
    )
