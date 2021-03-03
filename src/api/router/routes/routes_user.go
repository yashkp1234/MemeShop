package routes

import (
	"net/http"

	"github.com/yashkp1234/MemeShop.git/api/controller"
)

//UsersRoutes represents all user routes
var UsersRoutes = []Route{
	{
		URI:          "/users/{id}",
		Method:       http.MethodGet,
		Handler:      controller.GetUser,
		AuthRequired: false,
	},
	{
		URI:          "/users",
		Method:       http.MethodPost,
		Handler:      controller.CreateUser,
		AuthRequired: false,
	},
	{
		URI:          "/users/{id}",
		Method:       http.MethodPut,
		Handler:      controller.UpdateUser,
		AuthRequired: true,
	},
	{
		URI:          "/users/{id}",
		Method:       http.MethodDelete,
		Handler:      controller.DeleteUser,
		AuthRequired: true,
	},
}
