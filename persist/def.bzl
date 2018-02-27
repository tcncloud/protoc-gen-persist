load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_context",
    "GoLibrary",
    "go_library",
)
load(
    "@io_bazel_rules_go//proto:compiler.bzl",
    "GoProtoCompiler",
    "go_proto_compile",
)
load(
    "@io_bazel_rules_go//proto:def.bzl",
    "get_imports",
)

# TODO: Check if it's acceptable to use private code from bazel rules go
load(
    "@io_bazel_rules_go//go/private:rules/rule.bzl",
    "go_rule",
)

def get_outputs(ctx, direct_sources, file_names):
    outputs = []
    for src in direct_sources:
        base = src.basename[:-len(".proto")]

        if ( base.endswith("_p") == True ):
            for file_name in file_names:
                new_file_name = "%s~/%s/%s" % (ctx.attr.name, ctx.attr.importpath, file_name)
                new_file_name = new_file_name.format(basename = base)
                outputs += [ctx.actions.declare_file(new_file_name)]

    return outputs    

def get_proto_file_paths(proto_files):
    proto_paths = []
    for f in proto_files:
        proto_paths.append(f.path)
    
    return proto_paths

def combine_inputs(*inputs):
    combined = depset()
    for i in inputs:
        combined += i
    return combined

def _persist_lib_proto_compile_impl(ctx):
    importpath = ctx.attr.importpath
    proto = ctx.attr.proto.proto

    # Declare files and add to go_srcs
    file_names = [
        "persist_lib/{basename}_query_handlers.persist.go",
        "persist_lib/{basename}_queries.persist.go",
    ]
    outputs = get_outputs(ctx, proto.direct_sources, file_names)
    file_name = "%s~/%s/persist_lib/pkg_level_definitions.persist.go" % (ctx.attr.name, importpath)
    outputs += [ctx.actions.declare_file(file_name)]

    # Determine the outpath
    outpath = outputs[0].dirname[:-len(importpath + "/persist_lib")]

    # Add arguments
    args = ctx.actions.args()
    args.add([
        "--protoc", ctx.file._protoc,
        "--importpath", importpath,
        "--out_path", outpath,
        "--plugin", ctx.file._persist_plugin,
    ])

    # Set options and add to args
    options = ["import_path={}".format(importpath)]
    options += ["persist_root={}".format(importpath)]
    args.add(options, before_each = "--option")

    # Add descriptor set
    descriptor_sets = proto.transitive_descriptor_sets
    args.add(descriptor_sets, before_each = "--descriptor_set")

    # Add outputs
    args.add(outputs, before_each = "--expected")

    # Add proto paths
    proto_paths = get_proto_file_paths(proto.direct_sources)
    args.add(proto_paths)

    ctx.actions.run(
        inputs = combine_inputs(
            [ctx.file._persist_plugin],
            [ctx.file._protoc],
            [ctx.file._go_protoc],
            descriptor_sets
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

persist_lib_proto_compile = rule(
    attrs = {
        "importpath": attr.string(
            mandatory = True,
        ),
        "proto": attr.label(
            mandatory = True,
            allow_files = True,
            single_file = True,  #TODO: should this be single?
        ),
        "_persist_plugin": attr.label(
            default = Label("//:protoc-gen-persist"),
            executable = True,
            allow_files = True,
            single_file = True,
            cfg = "host",
        ),
        "_go_plugin": attr.label(
            allow_files = True,
            single_file = True,
            executable = True,
            cfg = "host",
            default = Label("@com_github_golang_protobuf//protoc-gen-go"),
        ),
        "_protoc": attr.label(
            allow_files = True,
            single_file = True,
            executable = True,
            cfg = "host",
            default = Label("@com_github_google_protobuf//:protoc"),
        ),
        "_go_protoc": attr.label(
            allow_files = True,
            single_file = True,
            executable = True,
            cfg = "host",
            default = Label("@io_bazel_rules_go//go/tools/builders:go-protoc"),
        ),
    },
    # output_to_genfiles = True, #Yes or no? attribute?
    implementation = _persist_lib_proto_compile_impl,
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
    persist_lib_proto_compile(
        name = persist_lib,
        importpath = importpath,
        proto = proto,
    )

    # Compile go library from persist lib files
    go_library(
        name = name,
        srcs = [":" + persist_lib],
        visibility = visibility,
        deps = go_lib_deps,
        importpath = importpath + "/persist_lib",
    )

#TODO: DRY
def _lib_proto_compile_impl(ctx):
    importpath = ctx.attr.importpath
    proto = ctx.attr.proto.proto

    # Declare files and add to outputs
    file_names = ["{basename}.persist.go"]
    outputs = get_outputs(ctx, proto.direct_sources, file_names)

    # Determine the outpath
    outpath = outputs[0].dirname[:-len(importpath)]

    # Add arguments
    args = ctx.actions.args()
    args.add([
        "--protoc", ctx.file._protoc,
        "--importpath", importpath,
        "--out_path", outpath,
        "--plugin", ctx.file._persist_plugin,
    ])

    # Set options and add to args
    options = ["import_path={}".format(importpath)]
    options += ["persist_root={}".format(importpath)]
    args.add(options, before_each = "--option")

    # Add descriptor set
    descriptor_sets = proto.transitive_descriptor_sets
    args.add(descriptor_sets, before_each = "--descriptor_set")

    # Add outputs
    args.add(outputs, before_each = "--expected")

    # Add proto paths
    proto_paths = get_proto_file_paths(proto.direct_sources)
    args.add(proto_paths)

    ctx.actions.run(
        inputs = combine_inputs(
            [ctx.file._persist_plugin],
            [ctx.file._protoc],
            [ctx.file._go_protoc],
            descriptor_sets
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

lib_proto_compile = rule(
    attrs = {
        "importpath": attr.string(
            mandatory = True,
        ),
        "proto": attr.label(
            mandatory = True,
            allow_files = True,
            single_file = True,  #TODO: should this be single?
        ),
        "_persist_plugin": attr.label(
            default = Label("//:protoc-gen-persist"),
            executable = True,
            allow_files = True,
            single_file = True,
            cfg = "host",
        ),
        "_go_plugin": attr.label(
            allow_files = True,
            single_file = True,
            executable = True,
            cfg = "host",
            default = Label("@com_github_golang_protobuf//protoc-gen-go"),
        ),
        "_protoc": attr.label(
            allow_files = True,
            single_file = True,
            executable = True,
            cfg = "host",
            default = Label("@com_github_google_protobuf//:protoc"),
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
    implementation = _lib_proto_compile_impl,
)

def _pb_go_compile_impl(ctx):
    go = go_context(ctx)
    compilers = ctx.attr.compilers
    go_srcs = []
    valid_archive = False
    for c in compilers:
        compiler = c[GoProtoCompiler]
        if compiler.valid_archive:
            valid_archive = True
        go_srcs.extend(compiler.compile(go,
            compiler = compiler,
            proto = ctx.attr.proto.proto,
            imports = get_imports(ctx.attr),
            importpath = go.importpath,
        ))

    return DefaultInfo(
        files=depset(go_srcs),
    )

pb_go_compile = go_rule(
    _pb_go_compile_impl,
    attrs = {
        "importpath": attr.string(
            mandatory = True,
        ),
        "deps": attr.label_list(
            providers = [GoLibrary],
            # aspects = [_go_proto_aspect],
        ),
        "proto": attr.label(
            mandatory = True,
            allow_files = True,
            single_file = True,
        ),
        "compilers": attr.label_list(
            providers = [GoProtoCompiler],
            default = ["@io_bazel_rules_go//proto:go_grpc"],
        ),
    },
    # output_to_genfiles = True,
)

def persist_go_library(
    name,
    go_srcs = [],
    go_lib_deps = [],
    importpath = None,
    visibility = None,
    proto = None,
    **kwargs):

    # Generate persist.go files in lib package
    persist_files = name + "_persist_files"
    lib_proto_compile(
        name = persist_files,
        importpath = importpath,
        proto = proto,
    )

    # Generate pb.go files in lib package
    pb_go_files = name + "_pb_files"
    pb_go_compile(
        name = pb_go_files,
        importpath = importpath,
        proto = proto,
    )

    # Compile go library from lib files
    go_library(
        name = name,
        srcs = [":" + persist_files] + [":" + pb_go_files] + go_srcs,
        importpath = importpath,
        visibility = visibility,
        deps = go_lib_deps,
    )
