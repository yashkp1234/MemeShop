package routes

import (
	"api/controller"
	"net/http"
)

//UsersRoutes represents all user routes
var UsersRoutes = []Route{
	{
		URI:     "/users",
		Method:  http.MethodGet,
		Handler: controller.GetUsers,
	},
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
