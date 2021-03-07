package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/yashkp1234/MemeShop.git/api/auth"
	"github.com/yashkp1234/MemeShop.git/api/models"
	"github.com/yashkp1234/MemeShop.git/api/responses"
)

//Login logs in a user
func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err := user.ValidateCreation(); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := auth.SignIn(user.UserName, user.Password)
	if err != nil {
		fmt.Println(err)
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("Invalid username or password"))
		return
	}

	responses.JSON(w, http.StatusOK, token)
}
