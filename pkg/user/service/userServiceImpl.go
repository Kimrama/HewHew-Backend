package service

import (
	"errors"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/user/model"
	"hewhew-backend/pkg/user/repository"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type UserServiceImpl struct {
	userRepository repository.UserRepository
}


func NewUserServiceImpl(userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

func (s *UserServiceImpl) EditUser(userID uuid.UUID, userEntity *entities.User) error {
	if userEntity.FName == "" && userEntity.LName == "" && userEntity.Gender == "" {
		return errors.New("no fields to update")
	}
	userEntity.UserID = userID
	return s.userRepository.EditUser(userID, userEntity)
}

func (s *UserServiceImpl) CreateUser(userModel *model.CreateUserRequest) error {
	imageUrl := "NULL"
	if userModel.Image != nil {
		var err error
		imageUrl, err = s.userRepository.UploadUserProfileImage(userModel.Username, userModel.Image)
		if err != nil {
			return err
		}
	}
	userEntity := &entities.User{
		UserID:          uuid.New(),
		Username:        userModel.Username,
		Password:        userModel.Password,
		FName:           userModel.FName,
		LName:           userModel.LName,
		Gender:          userModel.Gender,
		ProfileImageURL: imageUrl,
	}

	if err := s.userRepository.CreateUser(userEntity); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) CreateAdmin(userModel *model.CreateAdminRequest) error {

	AdminEntity := &entities.ShopAdmin{
		AdminID:         uuid.New(),
		Username:        userModel.Username,
		Password:        userModel.Password,
		FName:           userModel.FName,
		LName:           userModel.LName,
	}

	if err := s.userRepository.CreateAdmin(AdminEntity); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) EditUserProfileImage(userID uuid.UUID, imageModel *utils.ImageModel) error {
	err := s.userRepository.EditUserProfileImage(userID, imageModel)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) GetUserByUsername(username string) (*entities.User, error) {
	return s.userRepository.GetUserByUsername(username)
}
func (s *UserServiceImpl) GetUserByUserID(userID uuid.UUID) (*entities.User, error) {
	return s.userRepository.GetUserByUserID(userID)
}


func (s *UserServiceImpl) GetShopAdminByUsername(username string) (*entities.ShopAdmin, error) {
	return s.userRepository.GetShopAdminByUsername(username)
}
