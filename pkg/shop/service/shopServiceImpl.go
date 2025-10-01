package service

import (
	"errors"
	"fmt"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/shop/model"
	"hewhew-backend/pkg/shop/repository"
	"hewhew-backend/utils"

	"github.com/google/uuid"
)

type ShopServiceImpl struct {
	ShopRepository repository.ShopRepository
}

func NewShopServiceImpl(ShopRepository repository.ShopRepository) ShopService {
	return &ShopServiceImpl{
		ShopRepository: ShopRepository,
	}
}
func (s *ShopServiceImpl) CreateCanteen(canteenModel interface{}) error {
	return s.ShopRepository.CreateCanteen(canteenModel)
}

func (s *ShopServiceImpl) EditCanteen(canteenName string, canteenEntity *entities.Canteen) error {
	if canteenName == "" {
		return fmt.Errorf("canteen name is required")
	}
	return s.ShopRepository.EditCanteen(canteenName, canteenEntity)
}

func (s *ShopServiceImpl) DeleteCanteen(canteenName string) error {
	if canteenName == "" {
		return fmt.Errorf("canteen name is required")
	}
	return s.ShopRepository.DeleteCanteen(canteenName)
}

func (s *ShopServiceImpl) GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error) {
	return s.ShopRepository.GetShopByAdminID(adminID)
}

func (s *ShopServiceImpl) ChangeState(body model.ChangeState, shopID uuid.UUID) error {

	var state bool
	switch body.State {
	case "open":
		state = true
	case "close":
		state = false
	default:
		return fmt.Errorf("invalid state value: %s", body.State)
	}

	return s.ShopRepository.ChangeState(state, shopID)
}

func (s *ShopServiceImpl) EditShopImage(shopID uuid.UUID, imageModel *utils.ImageModel) error {
	err := s.ShopRepository.EditShopImage(shopID, imageModel)
	if err != nil {
		return err
	}
	return nil
}

func (s *ShopServiceImpl) EditShop(body model.EditShopRequest, shop uuid.UUID) error {
	if body.ShopName == "" && body.ShopCanteenName == "" {
		return errors.New("no fields to update")
	}

	shopEntity := &entities.Shop{
		Name:        body.ShopName,
		CanteenName: body.ShopCanteenName,
		Address:     "Null",
	}
	fmt.Println("Service - EditShop: ", shopEntity, shop)

	return s.ShopRepository.EditShop(*shopEntity, shop)
}

func (s *ShopServiceImpl) CreateTag(ShopID string, body *model.TagCreateRequest) (*entities.Tag, error) {
	if body.Topic == "" {
		return nil, errors.New("tag topic is required")
	}

	shopUUID, err := uuid.Parse(ShopID)
	if err != nil {
		return nil, fmt.Errorf("invalid ShopID: %v", err)
	}

	tagEntity := &entities.Tag{
		Topic:  body.Topic,
		ShopID: shopUUID,
		TagID:  uuid.New(),
	}

	return s.ShopRepository.CreateTag(tagEntity)
}

func (s *ShopServiceImpl) GetShopAdminByUsername(username string) (*entities.ShopAdmin, error) {
	return s.ShopRepository.GetShopAdminByUsername(username)
}

func (s *ShopServiceImpl) GetAllCanteens() ([]entities.Canteen, error) {
	return s.ShopRepository.GetAllCanteens()
}

func (s *ShopServiceImpl) GetTagsByShopIDAndTopic(shopID string, topic string) ([]entities.Tag, error) {
	return s.ShopRepository.GetTagsByShopIDAndTopic(shopID, topic)
}

func (s *ShopServiceImpl) EditTag(tagID string, topic string) error {
	tagUUID, err := uuid.Parse(tagID)
	if err != nil {
		return fmt.Errorf("invalid TagID: %v", err)
	}

	tagEntity := &entities.Tag{
		TagID: tagUUID,
		Topic: topic,
	}

	return s.ShopRepository.EditTag(tagEntity)
}

func (s *ShopServiceImpl) GetAllTags(shopID string) ([]entities.Tag, error) {
	return s.ShopRepository.GetAllTags(shopID)
}
func (s *ShopServiceImpl) DeleteTag(tagID string) error {
	err := s.ShopRepository.DeleteTag(tagID)
	if err != nil {
		return fmt.Errorf("failed to delete tag: %v", err)
	}
	return nil
}

func (s *ShopServiceImpl) GetAllMenus(shopID uuid.UUID) ([]*entities.Menu, error) {
	return s.ShopRepository.GetAllMenus(shopID)
}
