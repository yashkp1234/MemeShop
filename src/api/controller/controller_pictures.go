package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yashkp1234/MemeShop.git/api/database"
	"github.com/yashkp1234/MemeShop.git/api/models"
	"github.com/yashkp1234/MemeShop.git/api/repository"
	"github.com/yashkp1234/MemeShop.git/api/repository/crud"
	"github.com/yashkp1234/MemeShop.git/api/responses"
	"github.com/yashkp1234/MemeShop.git/api/utils/contextkey"
)

func getUsername(w http.ResponseWriter, r *http.Request) (string, error) {
	username, err := contextkey.GetCallerFromContext(r.Context())
	if !err {
		return "", errors.New("No username found in ctx")
	}
	return username, nil
}

// GetPicture lists a single picture
func GetPicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pictureID := vars["pid"]

	username, err := getUsername(w, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := database.Connect()
	repo := crud.NewRepositoryPicturesCRUD(db)

	func(picturesRepository repository.PictureRepository) {
		//Find picture
		picture, err := picturesRepository.FindByID(username, pictureID)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, picture.ID))
		responses.JSON(w, http.StatusCreated, picture)
	}(repo)
}

// CreatePicture creates a picture
func CreatePicture(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	picture := models.Picture{}
	err = json.Unmarshal(body, &picture)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := database.Connect()
	repo := crud.NewRepositoryPicturesCRUD(db)

	func(picturesRepository repository.PictureRepository) {
		picture, err := picturesRepository.Save(picture)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, picture.ID))
		responses.JSON(w, http.StatusCreated, picture)
	}(repo)
}

// UpdatePicture updates a picture
func UpdatePicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["pid"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var pictureUpdate map[string]string
	err = json.Unmarshal(body, &pictureUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if err = models.ValidatePictureUpdate(pictureUpdate); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	username, err := getUsername(w, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := database.Connect()
	repo := crud.NewRepositoryPicturesCRUD(db)

	func(picturesRepository repository.PictureRepository) {
		//Find picture
		err := picturesRepository.Update(id, username, pictureUpdate)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s/%s", r.Host, r.RequestURI, id))
		responses.JSON(w, http.StatusCreated, id)
	}(repo)
}

// DeletePicture deletes a picture
func DeletePicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["pid"]

	username, err := getUsername(w, r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := database.Connect()
	repo := crud.NewRepositoryPicturesCRUD(db)

	func(picturesRepository repository.PictureRepository) {
		//Find picture
		id, err := picturesRepository.Delete(username, id)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s/%s", r.Host, r.RequestURI, id))
		responses.JSON(w, http.StatusCreated, id)
	}(repo)
}
