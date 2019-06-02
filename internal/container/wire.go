//+build wireinject

package container

import (
	gopher "github.com/friendsofgo/gopherapi/pkg"
	"github.com/friendsofgo/gopherapi/pkg/adding"
	"github.com/friendsofgo/gopherapi/pkg/fetching"
	"github.com/friendsofgo/gopherapi/pkg/modifying"
	"github.com/friendsofgo/gopherapi/pkg/removing"
	"github.com/friendsofgo/gopherapi/pkg/server"
	"github.com/friendsofgo/gopherapi/pkg/storage/inmem"
	"github.com/google/wire"
)

func InitializeServer(gophers map[string]gopher.Gopher) server.Server {
	wire.Build(server.New, inmem.NewRepository, fetching.NewService, adding.NewService, modifying.NewService, removing.NewService)
	return server.NewWire()
}
