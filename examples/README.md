## Examples

This directory has examples of how to use protoc-gen-persist, what is generated, and the rules when constructing
the protobuf definitions.  There are 4 sub-directories in this directory.  3 of which are important:
- ___mytime___ this is an example go lib for custom type mapping of protobuf types
- ___sql___ contains examples of generated sql code, and example service implementations
- ___spanner___ container of our generated spanner code, and example service implementations. Spanner sql queries are
different than sql queries and have several limitations compared to normal sql queries.  This differences are explained in the examples

Each of the sub-directories in the sql, and spanner assume you ran the protoc commands:
- ``` protoc -I/usr/local/include -I. -I$GOPATH/src --go_out=plugins=grpc:. ./simple_service.proto ```
- ``` protoc -I/usr/local/include -I. -I$GOPATH/src --persist_out=plugins=protoc-gen-persist:. ./simple_service.proto ```

in each of the directories to generate the code from the proto files. The files are commented explaining what will
be generated from the proto files, but you should read the [documentation](../docs/index.md) for a more detailed walkthrough
