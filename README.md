# protoc-gen-persist

Protobuf plugin that generate a persistence layer backed by SQL or Spanner database


## Rationale
An opinionated protoc plugin that will help generate boiler plate go code for GRPC micro services projects that need to interact with a SQL or Spanner database.
The code generation is accomplished by providing a protobuf file with proper annotations, and then running the protoc cli tool with the `--go_out` option, and the `--persist_out` options set to the same directory.

## Installation
This project requires [protoc](https://github.com/google/protobuf) and [protoc-gen-go](https://developers.google.com/protocol-buffers/docs/gotutorial) be installed

Then install with ```go get github.com/tcncloud/protoc-gen-persist```
## Documentation
documentation for the project can be found [Here](docs/index.md)
The documentation goes over the persist options, how to structure your proto file,
custom type mapping to/from the database, and spanner query parsing

## Version 4.0.0
Starting with this version we changed the command line parameters.
On the `--plugin_out` option you can add `path=source_relative` to make sure that all the generated files will be created on the current directory.
The default behaviour (don't specify the `path=source_relative`) will generate the files in the same directory with the source files or in a directory computed from the  (persist.pkg) or go_package option.


## Version 3.0.0
- Complete rework of the plugin (check the [docs](docs/index.md), and [examples](https://github.com/tcncloud/protoc-gen-persist/tree/master/examples) directory for more info)

## Version 2.0.0
- new method for generating before and after hooks. See #82 (@iamneal)

## Version 1.0.0
- new persist_lib generated package for custom handlers
- generated service handlers and custom handlers can be on different types,
in different packages
- expanded default type mapping
- lots of bug fixes


## Authors
 * [Florin Stan](https://github.com/namtzigla)
 * [Neal Cooper](https://github.com/iamneal)
 * [Jamie Whahlin](https://github.com/jwahlin)
 * [Michael Sorenson](https://github.com/michael-the-grey)
 * [Colton Morris](https://github.com/coltonmorris)
 * Spencer Beard


## License
Copyright 2017, TCN Inc.
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

 * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
 * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
 * Neither the name of TCN Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
