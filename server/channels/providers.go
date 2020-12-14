package channels

import "github.com/google/wire"

var Providers = wire.NewSet(NewStore, HttpVmListHandler, HttpConnectDiskHandler)
