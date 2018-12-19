# protoc-gen-persist

[![Build Status](https://www.travis-ci.org/tcncloud/protoc-gen-persist.svg?branch=master)](https://www.travis-ci.org/tcncloud/protoc-gen-persist)

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

## NOTE: This project is under development, use it at your own risk

## Version 2.0.0
- new method for generating before and after hooks. See #82 (@iamneal)

## Version 1.0.0
- new persist_lib generated package for custom handlers
- generated service handlers and custom handlers can be on different types,
in different packages
- expanded default type mapping
- lots of bug fixes

## Roadmap
 1. ~~type mapping~~
 1. ~~before & after callback function~~
 1. add tests, lots of tests
 1. ~~rewrite hooks/type mappings system to be more import safe~~


## Authors
 * [Florin Stan](https://github.com/namtzigla)
 * [Neal Cooper](https://github.com/iamneal)
 * [Jamie Whahlin](https://github.com/jwahlin)


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
