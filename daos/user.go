package daos

import (
	"simpleRest/app"
	"simpleRest/models"
)

type UserDAO struct{}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

func (dao *UserDAO) Get(rs app.RequestScope, id int) (*models.User, error) {
	var user models.User
	err := rs.Tx().First(&user, id).Error
	return &user, err
}

func (dao *UserDAO) GetByEmail(rs app.RequestScope, email string) (*models.User, error) {
	var user models.User
	err := rs.Tx().First(&user, "business_email=?", email).Error
	return &user, err
}

func (dao *UserDAO) New(rs app.RequestScope, user *models.User) error {
	err := rs.Tx().Create(&user).Error
	return err
}

func (dao *UserDAO) Update(rs app.RequestScope, user *models.User) error {
	err := rs.Tx().Save(&user).Error
	return err
}
