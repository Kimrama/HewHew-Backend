package service

import (
	"errors"
	"fmt"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/menu/model"
	"hewhew-backend/pkg/menu/repository"
	"hewhew-backend/utils"
	"strconv"

	"github.com/google/uuid"
)

type MenuServiceImpl struct {
	MenuRepository repository.MenuRepository
}

func NewMenuServiceImpl(MenuRepository repository.MenuRepository) MenuService {
	return &MenuServiceImpl{
		MenuRepository: MenuRepository,
	}
}

func (s *MenuServiceImpl) CreateMenu(menuModel *model.MenuRequest, shopID uuid.UUID) error {
	menuID := uuid.New()
	imageUrl := ""
	if menuModel.Image != nil {
		var err error
		imageUrl, err = s.MenuRepository.UploadMenuImage(menuID, menuModel.Image)
		if err != nil {
			return err
		}
	}

	price, err := strconv.ParseFloat(menuModel.Price, 64)
	if err != nil {
		return errors.New("invalid price")
	}

	var tag1UUID, tag2UUID *uuid.UUID
	if menuModel.Tag1ID != "" {
		id, err := uuid.Parse(menuModel.Tag1ID)
		if err != nil {
			return err
		}
		tag1UUID = &id
	}
	if menuModel.Tag2ID != "" {
		id, err := uuid.Parse(menuModel.Tag2ID)
		if err != nil {
			return err
		}
		tag2UUID = &id
	}

	menuEntity := &entities.Menu{
		MenuID:   menuID,
		ShopID:   shopID,
		Name:     menuModel.Name,
		Detail:   menuModel.Detail,
		Price:    price,
		Status:   string(menuModel.Status),
		ImageURL: imageUrl,
		Tag1ID:   tag1UUID,
		Tag2ID:   tag2UUID,
	}

	if err := s.MenuRepository.CreateMenu(menuEntity); err != nil {
		return err
	}

	return nil
}

func (s *MenuServiceImpl) GetMenusByShopID(shopID uuid.UUID) ([]*entities.Menu, error) {
	return s.MenuRepository.GetMenusByShopID(shopID)
}

func (s *MenuServiceImpl) GetMenuByID(menuID uuid.UUID) (*entities.Menu, []string, error) {
	menu, err := s.MenuRepository.GetMenuByID(menuID)
	if err != nil {
		return nil, nil, err
	}

	var tags []string

	if menu.Tag1ID != nil {
		tag1, err := s.MenuRepository.GetTagByID(*menu.Tag1ID)
		if err == nil {
			tags = append(tags, tag1.Topic)
		}
	}

	if menu.Tag2ID != nil {
		if menu.Tag1ID == nil || *menu.Tag2ID != *menu.Tag1ID {
			tag2, err := s.MenuRepository.GetTagByID(*menu.Tag2ID)
			if err == nil {
				tags = append(tags, tag2.Topic)
			}
		}
	}

	return menu, tags, nil
}
func (s *MenuServiceImpl) DeleteMenu(menuID uuid.UUID, admin *entities.ShopAdmin) error {
	menu, err := s.MenuRepository.GetMenuByID(menuID)
	if err != nil {
		return fmt.Errorf("menu not found")
	}

	if menu.ShopID != admin.ShopID {
		return fmt.Errorf("unauthorized to delete this menu")
	}

	return s.MenuRepository.DeleteMenu(menuID)
}

func (s *MenuServiceImpl) EditMenu(menuID uuid.UUID, admin *entities.ShopAdmin, req *model.MenuRequest) error {
	price, err := strconv.ParseFloat(req.Price, 64)
	if err != nil {
		return errors.New("invalid price")
	}

	menu, err := s.MenuRepository.GetMenuByID(menuID)
	if err != nil {
		return err
	}
	if menu.ShopID != admin.ShopID {
		return errors.New("unauthorized to edit this menu")
	}

	var tag1UUID, tag2UUID *uuid.UUID
	if req.Tag1ID != "" {
		id, err := uuid.Parse(req.Tag1ID)
		if err != nil {
			return err
		}
		tag1UUID = &id
	}
	if req.Tag2ID != "" {
		id, err := uuid.Parse(req.Tag2ID)
		if err != nil {
			return err
		}
		tag2UUID = &id
	}

	menu = &entities.Menu{
		MenuID: menuID,
		Name:   req.Name,
		Detail: req.Detail,
		Price:  price,
		Tag1ID: tag1UUID,
		Tag2ID: tag2UUID,
	}

	return s.MenuRepository.EditMenu(menu)
}

func (s *MenuServiceImpl) EditMenuStatus(menuID uuid.UUID, admin *entities.ShopAdmin, status string) error {
	menu, err := s.MenuRepository.GetMenuByID(menuID)
	if err != nil {
		return err
	}
	if menu.ShopID != admin.ShopID {
		return errors.New("unauthorized to edit this menu")
	}
	return s.MenuRepository.EditMenuStatus(menuID, status)
}

func (s *MenuServiceImpl) EditMenuImage(menuID uuid.UUID, admin *entities.ShopAdmin, imageModel *utils.ImageModel) error {
	menu, err := s.MenuRepository.GetMenuByID(menuID)
	if err != nil {
		return err
	}
	if menu.ShopID != admin.ShopID {
		return errors.New("unauthorized to edit this menu")
	}
	err = s.MenuRepository.EditMenuImage(menuID, imageModel)
	if err != nil {
		return err
	}
	return nil
}
