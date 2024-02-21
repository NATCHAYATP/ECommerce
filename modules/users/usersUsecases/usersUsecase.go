package usersUsecases

import (
	"fmt"

	"github.com/NATCHAYATP/E-Commerce/config"
	"github.com/NATCHAYATP/E-Commerce/modules/users"
	"github.com/NATCHAYATP/E-Commerce/modules/users/usersRepositories"
	"github.com/NATCHAYATP/E-Commerce/pkg/ecommerceauth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Insert user
	result, err := u.usersRepository.InsertUser(req, false)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	// Find user
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}

	// Sign token
	accessToken, err := ecommerceauth.NewEcommerceAuth(ecommerceauth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
		RoleId: user.RoleId,
	})
	// Sign token
	refreshToken, err := ecommerceauth.NewEcommerceAuth(ecommerceauth.Refresh, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
		RoleId: user.RoleId,
	})

	// Set passport
	passport := &users.UserPassport{
		User: &users.User{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken: accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}
	return passport, nil
}
