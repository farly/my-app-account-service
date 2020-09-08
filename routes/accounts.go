package accounts

import (
	schema "accounts/datastore/schema"
	models "accounts/models"
	server "accounts/server"
	util "accounts/utils"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
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
	accountRouter.router.HandleFunc("/accounts", ListAccount).Methods(http.MethodGet)

	createAccountRoute := accountRouter.router.Methods(http.MethodPost).Subrouter()
	createAccountRoute.HandleFunc("/accounts", accountRouter.CreateAccount).Methods(http.MethodPost)
	createAccountRoute.Use(Validate)

}

func ListAccount(response http.ResponseWriter, request *http.Request) {
	fmt.Fprint(response, "List Accounts")
}

func (accountRouter *AccountRoute) CreateAccount(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(KeyUser{}).(schema.User)

	err := accountRouter.userModel.Create(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		util.ToJSON(server.Fail{Ok: false}, w)
	}

	w.WriteHeader(http.StatusOK)
	util.ToJSON(server.Success{Ok: true}, w)
}

func Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := &schema.User{}

		err := util.FromJSON(user, r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			util.ToJSON(server.Fail{Ok: false}, w)
		}

		// move this out so this would be resuable
		var validate = validator.New()
		// this registers the json fieldname
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		// end
		err = validate.Struct(user)

		if err != nil {
			formattedErrs := make(map[string]string)
			w.WriteHeader(http.StatusBadRequest)

			errs := err.(validator.ValidationErrors)

			for _, err := range errs {
				// cast the FieldError into our ValidationError and append to the slice
				ve := err.(validator.FieldError)
				// fmtErr := fmt.Sprintf(
				// 	"Field validation for '%s' failed on the '%s' tag",
				// 	ve.Field(),
				// 	ve.Tag(),
				// )
				// formattedErrs = append(formattedErrs, fmtErr)
				formattedErrs[ve.Field()] = ve.Tag()
			}

			util.ToJSON(server.Fail{
				Ok:     false,
				Errors: formattedErrs,
			}, w)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUser{}, *user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
