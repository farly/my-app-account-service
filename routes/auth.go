package accounts

import (
	auth "accounts/auth"
	models "accounts/datastore/models"
	server "accounts/server"
	util "accounts/utils"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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
		return
	}

	isAuthorized := user.IsAuthorized(auth.Password)

	if !isAuthorized {
		w.WriteHeader(http.StatusForbidden)
		util.ToJSON(server.Fail{Ok: false}, w)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
	})
	tokenString, err := token.SignedString([]byte("jalsdfjskdfjalfasjfalsfjaslf"))

	if err != nil {
		w.WriteHeader(http.StatusOK)
		util.ToJSON(server.AuthenticationSuccess{
			Ok:    isAuthorized,
			Token: err.Error(),
		}, w)
		return
		// log fatal error here
	}

	w.WriteHeader(http.StatusOK)
	util.ToJSON(server.AuthenticationSuccess{
		Ok:    isAuthorized,
		Token: tokenString,
	}, w)
}
