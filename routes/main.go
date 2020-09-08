package accounts

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type MainRoute struct {
	router *mux.Router
}

func NewMainRoute(router *mux.Router) *MainRoute {
	return &MainRoute{router}
}

func (mainRouter *MainRoute) SetupRoutes() {
	mainRouter.router.HandleFunc("/", SendResponse).Methods(http.MethodGet)
}

func SendResponse(response http.ResponseWriter, request *http.Request) {
	fmt.Fprint(response, "Root path")
}
