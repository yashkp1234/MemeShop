package routes

import (
	"net/http"

	"github.com/yashkp1234/MemeShop.git/api/controller"
)

//PicturesRoutes represents all picture routes
var PicturesRoutes = []Route{
	{
		URI:          "/pictures/{id}/{pid}",
		Method:       http.MethodGet,
		Handler:      controller.GetPicture,
		AuthRequired: true,
	},
	{
		URI:          "/pictures/{id}",
		Method:       http.MethodPost,
		Handler:      controller.CreatePicture,
		AuthRequired: true,
	},
	{
		URI:          "/pictures/{id}/{pid}",
		Method:       http.MethodPut,
		Handler:      controller.UpdatePicture,
		AuthRequired: true,
	},
	{
		URI:          "/pictures/{id}/{pid}",
		Method:       http.MethodDelete,
		Handler:      controller.DeletePicture,
		AuthRequired: true,
	},
}
