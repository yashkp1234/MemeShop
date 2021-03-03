package routes

import (
	"net/http"

	"github.com/yashkp1234/MemeShop.git/api/controller"
)

var loginRoutes = []Route{
	{
		URI:          "/login",
		Method:       http.MethodPost,
		Handler:      controller.Login,
		AuthRequired: false,
	},
}
