package repository

import (
	"bytes"
	"fmt"
	"hewhew-backend/config"
	"hewhew-backend/database"
	"hewhew-backend/entities"
	"hewhew-backend/utils"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopRepositoryImpl struct {
	db             database.Database
	supabaseConfig *config.Supabase
}

func NewShopRepositoryImpl(db database.Database, supabaseConfig *config.Supabase) ShopRepository {
	return &ShopRepositoryImpl{
		db:             db,
		supabaseConfig: supabaseConfig,
	}
}

func (r *ShopRepositoryImpl) CreateCanteen(canteenModel interface{}) error {
	return r.db.Connect().Create(canteenModel).Error
}

func (r *ShopRepositoryImpl) EditCanteen(canteenName string, canteen *entities.Canteen) error {
	db := r.db.Connect()
	err := db.Model(&entities.Canteen{}).
		Where("canteen_name = ?", canteenName).
		Updates(map[string]interface{}{
			"latitude":  canteen.Latitude,
			"longitude": canteen.Longitude,
		}).Error
	return err
}

func (r *ShopRepositoryImpl) DeleteCanteen(canteenName string) error {
	db := r.db.Connect()
	err := db.Where("canteen_name = ?", canteenName).Delete(&entities.Canteen{}).Error
	return err

}

func (r *ShopRepositoryImpl) GetCanteenByName(canteenName string) (*entities.Canteen, error) {
	var canteen entities.Canteen
	db := r.db.Connect()
	if err := db.Where("canteen_name = ?", canteenName).First(&canteen).Error; err != nil {
		return nil, err
	}
	return &canteen, nil
}

func (r *ShopRepositoryImpl) EditShop(body entities.Shop, shop uuid.UUID) error {
	db := r.db.Connect()
	err := db.Model(&entities.Shop{}).
		Where("shop_id = ?", shop).
		Updates(map[string]interface{}{
			"name":         body.Name,
			"canteen_name": body.CanteenName,
			"address":      body.Address,
		}).Error
	return err
}

func (r *ShopRepositoryImpl) ChangeState(state bool, shopID uuid.UUID) error {
	db := r.db.Connect()
	err := db.Model(&entities.Shop{}).
		Where("shop_id = ?", shopID).
		Update("state", state).Error
	return err
}

func (r *ShopRepositoryImpl) UploadShopImage(shopID uuid.UUID, imageModel *utils.ImageModel) (string, error) {
	customName := shopID.String() + "_" + fmt.Sprintf("%d", time.Now().Unix())

	mimeType := mime.TypeByExtension(imageModel.Ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	url := fmt.Sprintf("%s/storage/v1/object/images/shopProfile/%s", r.supabaseConfig.URL, customName)
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
		return "", fmt.Errorf("failed to upload image: %s", resp.Status)
	}
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/images/shopProfile/%s", r.supabaseConfig.URL, customName)
	return publicURL, nil
}

func (r *ShopRepositoryImpl) GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error) {
	// ดึง ShopID จากตาราง ShopAdmin
	var admin entities.ShopAdmin
	db := r.db.Connect()

	if err := db.Select("shop_id").First(&admin, "admin_id = ?", adminID).Error; err != nil {
		return nil, err
	}
	var shop entities.Shop
	if err := db.First(&shop, "shop_id = ?", admin.ShopID).Error; err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *ShopRepositoryImpl) EditShopImage(adminID uuid.UUID, imageModel *utils.ImageModel) error {
	db := r.db.Connect()

	fmt.Println("Admin ID in EditShopImage:", adminID)

	shop, err := r.GetShopByAdminID(adminID)
	if err != nil {
		return err
	}

	if shop.ImageURL != "" && shop.ImageURL != "NULL" {
		publicPrefixRender := fmt.Sprintf("%s/storage/v1/render/image/public/", r.supabaseConfig.URL)
		objectPath := strings.TrimPrefix(shop.ImageURL, publicPrefixRender)

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

	newImageURL, err := r.UploadShopImage(shop.ShopID, imageModel) // <- (string, error)
	if err != nil {
		return err
	}

	if err := db.Model(&entities.Shop{}).
		Where("shop_id = ?", shop.ShopID).
		Update("image_url", newImageURL).Error; err != nil {
		return err
	}

	return nil
}

func (r *ShopRepositoryImpl) GetShopAdminByUsername(username string) (*entities.ShopAdmin, error) {
	var admin entities.ShopAdmin
	db := r.db.Connect()
	if err := db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *ShopRepositoryImpl) CreateTag(tagModel *entities.Tag) (*entities.Tag, error) {
	if err := r.db.Connect().Create(tagModel).Error; err != nil {
		return nil, err
	}
	return tagModel, nil
}

func (r *ShopRepositoryImpl) CreateTransactionLog(log *entities.TransactionLog) error {
	return r.db.Connect().Create(log).Error
}

func (r *ShopRepositoryImpl) GetAllCanteens() ([]entities.Canteen, error) {
	var canteens []entities.Canteen
	db := r.db.Connect()
	if err := db.Find(&canteens).Error; err != nil {
		return nil, err
	}
	return canteens, nil
}

func (r *ShopRepositoryImpl) GetAllShops() ([]entities.Shop, error) {
	var shops []entities.Shop
	db := r.db.Connect()

	if err := db.Preload("Tags").Find(&shops).Error; err != nil {
		return nil, err
	}

	return shops, nil
}

func (r *ShopRepositoryImpl) GetNearbyShops(latitude, longitude float64) ([]entities.Shop, error) {
	var shops []entities.Shop
	db := r.db.Connect()
	if err := db.Where("state = ?", true).Find(&shops).Error; err != nil {
		return nil, err
	}
	return shops, nil
}

func (r *ShopRepositoryImpl) GetShopsByCanteens(canteenNames []string) ([]entities.Shop, error) {
	var shops []entities.Shop
	db := r.db.Connect()
	if err := db.Where("canteen_name IN ?", canteenNames).Find(&shops).Error; err != nil {
		return nil, err
	}
	return shops, nil
}

func (r *ShopRepositoryImpl) GetShopByID(shopID uuid.UUID) (*entities.Shop, error) {
	db := r.db.Connect()
	var shop entities.Shop
	err := db.Preload("Menus").Preload("Tags").First(&shop, "shop_id = ?", shopID).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *ShopRepositoryImpl) GetAllMenus(shopID uuid.UUID) ([]*entities.Menu, error) {
	var menus []*entities.Menu
	db := r.db.Connect()
	if err := db.Where("shop_id = ?", shopID).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *ShopRepositoryImpl) GetTagsByShopIDAndTopic(shopID string, topic string) ([]entities.Tag, error) {
	var tags []entities.Tag
	db := r.db.Connect()
	if err := db.
		Where("shop_id = ? AND topic = ?", shopID, topic).
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *ShopRepositoryImpl) EditTag(tagModel *entities.Tag) error {
	db := r.db.Connect()
	err := db.Model(&entities.Tag{}).
		Where("tag_id = ?", tagModel.TagID).
		Updates(map[string]interface{}{
			"topic": tagModel.Topic,
		}).Error
	return err
}

func (r *ShopRepositoryImpl) GetAllTags(shopID uuid.UUID) ([]entities.Tag, error) {
	var tags []entities.Tag
	db := r.db.Connect()
	if err := db.Where("shop_id = ?", shopID).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
func (r *ShopRepositoryImpl) DeleteTag(tagID string) error {
	db := r.db.Connect()
	err := db.Where("tag_id = ?", tagID).Delete(&entities.Tag{}).Error
	return err
}

func (r *ShopRepositoryImpl) CreateNotification(notification *entities.Notification) error {
	return r.db.Connect().Create(notification).Error
}

func (r *ShopRepositoryImpl) GetDropOffByID(id uuid.UUID) (*entities.DropOffLocation, error) {
	db := r.db.Connect()
	var dropOff entities.DropOffLocation
	err := db.First(&dropOff, "drop_off_location_id = ?", id).Error
	return &dropOff, err
}

func (r *ShopRepositoryImpl) GetOrderByID(orderID uuid.UUID) (*entities.Order, error) {
	db := r.db.Connect()
	var order entities.Order
	err := db.Preload("MenuQuantity").First(&order, "order_id = ?", orderID).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *ShopRepositoryImpl) GetMenuByID(menuID uuid.UUID) (*entities.Menu, error) {
	db := r.db.Connect()
	var menu entities.Menu
	err := db.First(&menu, "menu_id = ?", menuID).Error
	return &menu, err
}

func (r *ShopRepositoryImpl) GetTagByID(tagID uuid.UUID) (*entities.Tag, error) {
	var tag entities.Tag
	db := r.db.Connect()
	if err := db.First(&tag, "tag_id = ?", tagID).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *ShopRepositoryImpl) GetOrderIDsFromTransactionLog() ([]uuid.UUID, error) {
	db := r.db.Connect()
	var orderIDs []uuid.UUID

	if err := db.Model(&entities.TransactionLog{}).
		Pluck("DISTINCT order_id", &orderIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch order IDs from transaction logs: %v", err)
	}

	return orderIDs, nil
}

func (r *ShopRepositoryImpl) CountMenusFromOrders(orderIDs []uuid.UUID) (map[uuid.UUID]int, error) {
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
func (r *ShopRepositoryImpl) GetPopularShopsByMenuCounts(menuCounts map[uuid.UUID]int) ([]*entities.Shop, error) {
	if len(menuCounts) == 0 {
		return []*entities.Shop{}, nil
	}

	db := r.db.Connect()
	var shops []*entities.Shop
	menuIDs := make([]uuid.UUID, 0, len(menuCounts))
	for id := range menuCounts {
		menuIDs = append(menuIDs, id)
	}

	err := db.Joins("JOIN menus ON menus.shop_id = shops.shop_id").
		Where("menus.menu_id IN ?", menuIDs).
		Group("shops.shop_id").
		Order(gorm.Expr("SUM(CASE WHEN menus.menu_id IN ? THEN menus.price ELSE 0 END) DESC", menuIDs)).
		Find(&shops).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch popular shops: %v", err)
	}

	return shops, nil
}