package routes

import (
	"net/http"

	"github.com/yashkp1234/MemeShop.git/api/controller"
)

//UsersRoutes represents all user routes
var UsersRoutes = []Route{
	{
		URI:     "/users/{id}",
		Method:  http.MethodGet,
		Handler: controller.GetUser,
	},
	{
		URI:     "/users",
		Method:  http.MethodPost,
		Handler: controller.CreateUser,
	},
	{
		URI:     "/users/{id}",
		Method:  http.MethodPut,
		Handler: controller.UpdateUser,
	},
	{
		URI:     "/users/{id}",
		Method:  http.MethodDelete,
		Handler: controller.DeleteUser,
	},
}
