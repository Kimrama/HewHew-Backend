package service

import (
	"errors"
	"fmt"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/user/model"
	"hewhew-backend/pkg/user/repository"
	"hewhew-backend/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	contactEntity := &entities.Contact{
		ContactID:   uuid.New(),
		UserID:      userEntity.UserID,
		ContactType: userModel.ContactType,
		Detail:      userModel.ContactDetail,
	}
	if err := s.userRepository.CreateContact(contactEntity); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) CreateAdmin(userModel *model.CreateAdminRequest) error {

	ShopEntity := &entities.Shop{
		ShopID:      uuid.New(),
		Name:        "Null",
		Address:     "Null",
		CanteenName: "Null",
	}

	if err := s.userRepository.CreateShop(ShopEntity); err != nil {
		return err
	}

	AdminEntity := &entities.ShopAdmin{
		AdminID:  uuid.New(),
		Username: userModel.Username,
		Password: userModel.Password,
		FName:    userModel.FName,
		LName:    userModel.LName,
		ShopID:   ShopEntity.ShopID,
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

func (s *UserServiceImpl) CreateUserContact(userID uuid.UUID, req *model.EditUserContactRequest) error {
	newContact := &entities.Contact{
		ContactID:   uuid.New(),
		UserID:      userID,
		ContactType: req.ContactType,
		Detail:      req.Detail,
	}

	if err := s.userRepository.CreateContact(newContact); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) DeleteUserContact(userID, contactID uuid.UUID) error {
	contact, err := s.userRepository.GetContactByID(contactID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("contact not found")
		}
		return err
	}

	if contact.UserID != userID {
		return fmt.Errorf("unauthorized: cannot delete another user's contact")
	}

	return s.userRepository.DeleteContact(contactID)
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

func (s *UserServiceImpl) GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error) {
	return s.userRepository.GetShopByAdminID(adminID)
}

func (s *UserServiceImpl) GetAllShops() ([]entities.Shop, error) {
	return s.userRepository.GetAllShops()
}

func (s *UserServiceImpl) Topup(UserID string, amount float64) error {
	parsedUUID, err := uuid.Parse(UserID)
	if err != nil {
		return fmt.Errorf("invalid UserID: %v", err)
	}
	TopupEntity := &entities.TopUp{
		UserID:    parsedUUID,
		Amount:    amount,
		TimeStamp: time.Now(),
		TopUpID:   uuid.New(),
	}

	return s.userRepository.Topup(TopupEntity)
}
