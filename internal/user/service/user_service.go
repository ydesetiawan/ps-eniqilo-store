package service

import (
	"ps-eniqilo-store/internal/user/dto"
	"ps-eniqilo-store/internal/user/model"
	"ps-eniqilo-store/internal/user/repository"
	"ps-eniqilo-store/pkg/bcrypt"
	"ps-eniqilo-store/pkg/errs"
	"ps-eniqilo-store/pkg/middleware"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserServiceImpl(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) RegisterUser(req dto.RegisterReq) (*dto.RegisterResp, error) {
	phoneNumber := req.PhoneNumber
	hashedPassword, _ := bcrypt.HashPassword(req.Password)
	req.Password = hashedPassword
	id, err := s.userRepository.RegisterUser(model.NewUser(req))

	if err != nil {
		return &dto.RegisterResp{}, err
	}
	token, _ := middleware.GenerateJWT(phoneNumber, id)
	return &dto.RegisterResp{
		PhoneNumber: req.PhoneNumber,
		Name:        req.Name,
		AccessToken: token,
	}, nil
}

func (s *UserService) Login(req dto.LoginReq) (*dto.RegisterResp, error) {
	//TODO validation request
	usr, err := s.userRepository.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return &dto.RegisterResp{}, errs.NewErrDataNotFound("user not found ", req.PhoneNumber, errs.ErrorData{})
	}
	err = bcrypt.ComparePassword(req.Password, usr.Password)
	if err != nil {
		return &dto.RegisterResp{}, errs.NewErrBadRequest("password is wrong ")
	}

	token, _ := middleware.GenerateJWT(usr.PhoneNumber, usr.ID)

	return &dto.RegisterResp{
		PhoneNumber: usr.PhoneNumber,
		Name:        usr.Name,
		AccessToken: token,
	}, nil
}
