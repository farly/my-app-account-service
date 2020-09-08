package accounts

import (
	server "accounts/server"

	"github.com/gorilla/mux"
)

func SetUpRoutes(router *mux.Router, context *server.Context) {
	var mainRoute = NewMainRoute(router)
	mainRoute.SetupRoutes()

	var accountsRoute = NewAccountRoute(router, context.UserModel)
	accountsRoute.SetupRoutes()
}
