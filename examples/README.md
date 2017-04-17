## Examples

This directory has examples of how to use protoc-gen-persist, what is generated, and the rules when constructing
the protobuf definitions.  There are 4 sub-directories in this directory.  3 of which are important:
- ___mytime___ this is an example go lib for custom type mapping of protobuf types
- ___sql___ contains examples of generated sql code, and example service implementations
- ___spanner___ container of our generated spanner code, and example service implementations. Spanner sql queries are
different than sql queries and have several limitations compared to normal sql queries.  This differences are explained in the examples

Each of the sub-directories in the sql, and spanner assume you ran the protoc commands to generate the code for the proto files in each
directory.  The files are thoroughly commented explaining what will be generated  on the proto files, and readmes in each of those directories
will point at the specifics.
