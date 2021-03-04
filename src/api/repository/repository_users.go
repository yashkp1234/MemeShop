package repository

import "github.com/yashkp1234/MemeShop.git/api/models"

//UserRepository represents all operations we can perform
//with a user object
type UserRepository interface {
	Save(models.User) (models.User, error)
	FindByID(string) (models.User, error)
	Update(string, bool, models.User) error
	Delete(string) (string, error)
}
