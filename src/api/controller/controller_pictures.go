package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yashkp1234/MemeShop.git/api/cache"
	"github.com/yashkp1234/MemeShop.git/api/database"
	"github.com/yashkp1234/MemeShop.git/api/gcp"
	"github.com/yashkp1234/MemeShop.git/api/models"
	"github.com/yashkp1234/MemeShop.git/api/repository"
	"github.com/yashkp1234/MemeShop.git/api/repository/crud"
	"github.com/yashkp1234/MemeShop.git/api/responses"
	"github.com/yashkp1234/MemeShop.git/api/utils/contextkey"
)

func readPictureFromReq(r *http.Request) (*multipart.File, *multipart.FileHeader, error) {
	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("file")
	if err != nil {
		return nil, nil, err
	}
	return &file, handler, nil
}

// GetPicture lists a single picture
func GetPicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pictureID := vars["pid"]

	username, _ := contextkey.GetUsernameFromContext(r)

	db := database.Connect()
	cache := cache.Connect()
	imgCloud := gcp.Connect()
	repo := crud.NewRepositoryPicturesCRUD(db, cache, imgCloud)

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
	r.ParseMultipartForm(32 << 20)

	body := r.MultipartForm.Value["request"][0]

	picture := models.Picture{}
	err := json.Unmarshal([]byte(body), &picture)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	picture.User, err = contextkey.GetUsernameFromContext(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	file, handler, err := readPictureFromReq(r)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	db := database.Connect()
	cache := cache.Connect()
	imgCloud := gcp.Connect()
	repo := crud.NewRepositoryPicturesCRUD(db, cache, imgCloud)

	func(picturesRepository repository.PictureRepository) {
		picture, err := picturesRepository.Save(picture, file, handler)
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

	username, err := contextkey.GetUsernameFromContext(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := database.Connect()
	cache := cache.Connect()
	imgCloud := gcp.Connect()
	repo := crud.NewRepositoryPicturesCRUD(db, cache, imgCloud)

	func(picturesRepository repository.PictureRepository) {
		//Update picture
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

	username, err := contextkey.GetUsernameFromContext(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	db := database.Connect()
	cache := cache.Connect()
	imgCloud := gcp.Connect()
	repo := crud.NewRepositoryPicturesCRUD(db, cache, imgCloud)

	func(picturesRepository repository.PictureRepository) {
		//Delete picture
		id, err := picturesRepository.Delete(username, id)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}
		w.Header().Set("Location", fmt.Sprintf("%s%s/%s", r.Host, r.RequestURI, id))
		responses.JSON(w, http.StatusCreated, id)
	}(repo)
}

//GetUserPictures gets all the pictures for a user
func GetUserPictures(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	searchUsername := vars["username"]

	authUsername, _ := contextkey.GetUsernameFromContext(r)

	db := database.Connect()
	cache := cache.Connect()
	imgCloud := gcp.Connect()
	repo := crud.NewRepositoryPicturesCRUD(db, cache, imgCloud)

	func(picturesRepository repository.PictureRepository) {
		//Find pictures
		pictures, err := picturesRepository.FindByUser(searchUsername, searchUsername == authUsername)
		if err != nil {
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("%s%s", r.Host, r.RequestURI))
		responses.JSON(w, http.StatusCreated, pictures)
	}(repo)

}
