# protoc-gen-persist

Protobuf plugin that generate a persistence layer backed by SQL database

## NOTE: This project is under development everything could be changed !

## Rationale
Opinionated protoc plugin that will help generate builerplate go code for grpc microservices that need to interact 
with a SQL or Spanner database.


protoc-gen-persist is a protoc plugin used to generate golang CRUD services that communicate through GRPC,                     and talk to a persistence layer back end.  This is accomplished by providing a protobuf file with proper                       
annotations, and then running the protoc cli tool with the ```--go_out``` option, and the ```--persist_out```                 
options set to the same directory.

## Installation

## Documentation

## Versions

## Roadmap
 1. ~~type mapping~~
 1. before & after callback function
 1. add tests, lots of tests
 1. mongo 


## Authors
 * Florin Stan (@namtzigla)
 * Neal Cooper (@iamneal)


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
