package basic

import "github.com/tcncloud/protoc-gen-persist/examples/spanner/basic/persist_lib"

type MySpannerImpl struct {
	PERSIST *persist_lib.MySpannerPersistHelper
}
