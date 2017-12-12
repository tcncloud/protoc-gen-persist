package bob_example

import "github.com/tcncloud/protoc-gen-persist/examples/spanner/bob_example/persist_lib"

type BobsImpl struct {
	PERSIST *persist_lib.BobsPersistHelper
}
