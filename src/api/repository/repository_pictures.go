package repository

import (
	"mime/multipart"

	"github.com/yashkp1234/MemeShop.git/api/models"
)

//PictureRepository represents all operations we can perform
//with a user object
type PictureRepository interface {
	Save(models.Picture, *multipart.File, *multipart.FileHeader) (models.Picture, error)
	FindByID(string, string) (models.Picture, error)
	//FindByUser(string) ([]models.Picture, error)
	Update(string, string, map[string]string) error
	Delete(string, string) (string, error)
}
