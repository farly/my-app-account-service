package accounts

import (
	models "accounts/datastore/models"
	schema "accounts/datastore/schema"
	server "accounts/server"
	util "accounts/utils"
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type AccountRoute struct {
	router    *mux.Router
	userModel *models.UserModel
}

type KeyUser struct{}

func NewAccountRoute(router *mux.Router, userModel *models.UserModel) *AccountRoute {
	return &AccountRoute{router, userModel}
}

func (accountRouter *AccountRoute) SetupRoutes() {
	createAccountRoute := accountRouter.router.Methods(http.MethodPost).Subrouter()
	createAccountRoute.HandleFunc("/accounts", accountRouter.CreateAccount).Methods(http.MethodPost)
	createAccountRoute.Use(Validate)

}

func (accountRouter *AccountRoute) CreateAccount(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(KeyUser{}).(schema.User)

	err := accountRouter.userModel.Create(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		util.ToJSON(server.Fail{Ok: false}, w)
		return
	}

	// payload user.. generate token

	w.WriteHeader(http.StatusOK)
	util.ToJSON(server.Success{Ok: true}, w)
}

func Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := &schema.User{}

		err := util.FromJSON(user, r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			util.ToJSON(server.Fail{
				Ok: false,
				Errors: map[string]string{
					"body": "required",
				},
			}, w)
			return
		}

		var validate = util.NewValidate()

		errors := validate.Validate(user)

		if len(errors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			util.ToJSON(server.Fail{
				Ok:     false,
				Errors: errors,
			}, w)
			return
		}
		ctx := context.WithValue(r.Context(), KeyUser{}, *user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
