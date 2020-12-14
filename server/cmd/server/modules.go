//+build wireinject

package main

import (
	"github.com/burbokop/architecture-lab3/server/channels"
	"github.com/google/wire"
)

// ComposeApiServer will create an instance of CharApiServer according to providers defined in this file.
func ComposeApiServer(port HttpPortNumber) (*ChatApiServer, error) {
	wire.Build(
		// DB connection provider (defined in main.go).
		NewDbConnection,
		// Add providers from channels package.
		channels.Providers,
		// Provide ChatApiServer instantiating the structure and injecting channels handler and port number.
		wire.Struct(new(ChatApiServer), "Port", "ListVmsHandler", "DiscConnectionHandler"),
	)
	return nil, nil
}
