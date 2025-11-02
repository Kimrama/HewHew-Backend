package repository

import (
	"bytes"
	"fmt"
	"hewhew-backend/config"
	"hewhew-backend/database"
	"hewhew-backend/entities"
	"hewhew-backend/utils"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MenuRepositoryImpl struct {
	db             database.Database
	supabaseConfig *config.Supabase
}

func NewMenuRepositoryImpl(db database.Database, supabaseConfig *config.Supabase) MenuRepository {
	return &MenuRepositoryImpl{
		db:             db,
		supabaseConfig: supabaseConfig,
	}
}

func (r *MenuRepositoryImpl) CreateMenu(menuEntity *entities.Menu) error {
	db := r.db.Connect()
	err := db.Create(menuEntity).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *MenuRepositoryImpl) UploadMenuImage(menuID uuid.UUID, imageModel *utils.ImageModel) (string, error) {
	customName := menuID.String() + "_" + fmt.Sprintf("%d", time.Now().Unix())

	mimeType := mime.TypeByExtension(imageModel.Ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	url := fmt.Sprintf("%s/storage/v1/object/images/menuImage/%s", r.supabaseConfig.URL, customName)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(imageModel.Body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.supabaseConfig.Key))
	req.Header.Set("Content-Type", mimeType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload image: %s, %s", resp.Status, string(body))
	}
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/images/menuImage/%s", r.supabaseConfig.URL, customName)
	return publicURL, nil
}

func (r *MenuRepositoryImpl) GetMenusByShopID(shopID uuid.UUID) ([]*entities.Menu, error) {
	var menus []*entities.Menu
	db := r.db.Connect()
	if err := db.Where("shop_id = ?", shopID).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *MenuRepositoryImpl) GetMenuByID(menuID uuid.UUID) (*entities.Menu, error) {
	var menu entities.Menu
	db := r.db.Connect()
	if err := db.Where("menu_id = ?", menuID).First(&menu).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *MenuRepositoryImpl) GetTagByID(tagID uuid.UUID) (*entities.Tag, error) {
	var tag entities.Tag
	db := r.db.Connect()
	if err := db.First(&tag, "tag_id = ?", tagID).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *MenuRepositoryImpl) DeleteMenu(menuID uuid.UUID) error {
	db := r.db.Connect()
	if err := db.Where("menu_id = ?", menuID).Delete(&entities.Menu{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *MenuRepositoryImpl) EditMenu(menu *entities.Menu) error {
	db := r.db.Connect()
	err := db.Model(&entities.Menu{}).
		Where("menu_id = ?", menu.MenuID).
		Updates(map[string]interface{}{
			"name":    menu.Name,
			"detail":  menu.Detail,
			"price":   menu.Price,
			"tag1_id": menu.Tag1ID,
			"tag2_id": menu.Tag2ID,
		}).Error
	return err
}

func (r *MenuRepositoryImpl) EditMenuStatus(menuID uuid.UUID, status string) error {
	db := r.db.Connect()
	err := db.Model(&entities.Menu{}).
		Where("menu_id = ?", menuID).
		Update("status", status).Error
	return err
}

func (r *MenuRepositoryImpl) EditMenuImage(menuID uuid.UUID, imageModel *utils.ImageModel) error {
	db := r.db.Connect()
	menu, err := r.GetMenuByID(menuID)
	if err != nil {
		return err
	}

	if menu.ImageURL != "NULL" && menu.ImageURL != "" {

		publicPrefixRender := fmt.Sprintf("%s/storage/v1/render/image/public/", r.supabaseConfig.URL)
		objectPath := strings.TrimPrefix(menu.ImageURL, publicPrefixRender)

		deleteURL := fmt.Sprintf("%s/storage/v1/object/%s", r.supabaseConfig.URL, objectPath)

		req, _ := http.NewRequest("DELETE", deleteURL, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.supabaseConfig.Key))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("failed to delete image: %s", resp.Status)
		}
	}

	newImageUrl, err := r.UploadMenuImage(menuID, imageModel)
	if err != nil {
		return err
	}

	err = db.Model(&entities.Menu{}).
		Where("menu_id = ?", menuID).
		Update("image_url", newImageUrl).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *MenuRepositoryImpl) GetOrderIDsFromTransactionLog() ([]uuid.UUID, error) {
	db := r.db.Connect()
	var orderIDs []uuid.UUID

	if err := db.Model(&entities.TransactionLog{}).
		Pluck("DISTINCT order_id", &orderIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch order IDs from transaction logs: %v", err)
	}

	return orderIDs, nil
}

func (r *MenuRepositoryImpl) CountMenusFromOrders(orderIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	if len(orderIDs) == 0 {
		return map[uuid.UUID]int{}, nil
	}

	db := r.db.Connect()
	type Result struct {
		MenuID   uuid.UUID
		TotalQty int
	}
	var results []Result

	if err := db.Table("menu_quantities").
		Select("menu_id, SUM(quantity) as total_qty").
		Where("order_id IN ?", orderIDs).
		Group("menu_id").
		Order("total_qty DESC").
		Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to count popular menus: %v", err)
	}

	menuCounts := make(map[uuid.UUID]int)
	for _, res := range results {
		menuCounts[res.MenuID] = res.TotalQty
	}

	return menuCounts, nil
}

func (r *MenuRepositoryImpl) GetMenusByIDs(menuCounts map[uuid.UUID]int) ([]*entities.Menu, error) {
	if len(menuCounts) == 0 {
		return []*entities.Menu{}, nil
	}

	db := r.db.Connect()
	var menus []*entities.Menu
	menuIDs := make([]uuid.UUID, 0, len(menuCounts))
	for id := range menuCounts {
		menuIDs = append(menuIDs, id)
	}

	if err := db.Where("menu_id IN ?", menuIDs).Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch menus by IDs: %v", err)
	}

	return menus, nil
}