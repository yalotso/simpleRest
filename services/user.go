package services

import (
	"errors"
	"simpleRest/app"
	"simpleRest/models"
	"simpleRest/utils"
)

type userDAO interface {
	Get(rs app.RequestScope, id int) (*models.User, error)
	GetByEmail(rs app.RequestScope, email string) (*models.User, error)
	New(rs app.RequestScope, user *models.User) error
	Update(rs app.RequestScope, user *models.User) error
}

type UserService struct {
	dao userDAO
}

func NewUserService(dao userDAO) *UserService {
	return &UserService{dao}
}

func (s *UserService) New(rs app.RequestScope, user *models.User) error {
	err := user.Validate()
	if err != nil {
		return err
	}
	user.EmailCode = utils.GenerateVerificationCode(user.BusinessEmail)
	err = s.dao.New(rs, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Verify(rs app.RequestScope, code string, userId int) error {
	user, err := s.dao.Get(rs, userId)
	if err != nil {
		return err
	}
	if user.EmailCode != code {
		return errors.New("wrong verification code")
	}
	user.IsVerified = true
	err = s.dao.Update(rs, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ForgotPassword(rs app.RequestScope, email string) (*models.User, error) {
	user, err := s.dao.GetByEmail(rs, email)
	if err != nil {
		return nil, err
	}
	user.PasswordCode = utils.GenerateVerificationCode(user.PasswordCode)
	err = s.dao.Update(rs, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) ResetPassword(rs app.RequestScope, code, pass, confirmPass string, userId int) error {
	user, err := s.dao.Get(rs, userId)
	if err != nil {
		return err
	}
	if code != user.PasswordCode {
		return errors.New("wrong password code")
	}
	user.Password = pass
	user.ConfirmPassword = confirmPass
	user.PasswordCode = ""
	err = user.Validate()
	if err != nil {
		return err
	}
	err = s.dao.Update(rs, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) Login(rs app.RequestScope, email, password string) (*models.User, error) {
	user, err := s.dao.GetByEmail(rs, email)
	if err != nil {
		return nil, err
	}
	if password != user.Password {
		return nil, errors.New("wrong password")
	}
	return user, nil
}
