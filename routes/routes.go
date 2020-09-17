package accounts

import (
	server "accounts/server"

	"github.com/gorilla/mux"
)

func SetUpRoutes(router *mux.Router, context *server.Context) {

	var accountsRoute = NewAccountRoute(router, context.UserModel)
	accountsRoute.SetupRoutes()

	var loginRoute = NewLoginRoute(router, context.UserModel)
	loginRoute.SetupRoutes()
}
