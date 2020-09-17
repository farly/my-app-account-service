package accounts

import (
	auth "accounts/auth"
	models "accounts/datastore/models"
	server "accounts/server"
	util "accounts/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type LoginRoute struct {
	router    *mux.Router
	userModel *models.UserModel
}

func NewLoginRoute(router *mux.Router, userModel *models.UserModel) *LoginRoute {
	return &LoginRoute{router, userModel}
}

func (loginRouter *LoginRoute) SetupRoutes() {
	loginRoute := loginRouter.router.Methods(http.MethodPost).Subrouter()
	loginRoute.HandleFunc("/login", loginRouter.Login).Methods(http.MethodPost)
	// createAccountRoute.Use(Validate)
}

func (loginRouter *LoginRoute) Login(w http.ResponseWriter, r *http.Request) {

	auth := &auth.Auth{}
	err := util.FromJSON(auth, r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		util.ToJSON(server.Fail{
			Ok: false,
			Errors: map[string]string{
				"error": "marshal error",
			},
		}, w)
		return
	}

	user, err := loginRouter.userModel.FindOneByUsername(auth.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		util.ToJSON(server.Fail{
			Ok: false,
		}, w)
	}

	isAuthorized := user.IsAuthorized(auth.Password)

	w.WriteHeader(http.StatusOK)
	util.ToJSON(server.Success{Ok: isAuthorized}, w)
}
