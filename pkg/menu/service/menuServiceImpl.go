package service

import (
	"fmt"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/menu/model"
	"hewhew-backend/pkg/menu/repository"
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

func (s *MenuServiceImpl) CreateMenu(menuModel *model.CreateMenuRequest, shopID uuid.UUID) error {
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
        return fmt.Errorf("invalid price: %v", err)
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

