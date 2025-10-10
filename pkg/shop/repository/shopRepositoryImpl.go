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
	publicURL := fmt.Sprintf("%s/storage/v1/render/image/public/images/shopProfile/%s", r.supabaseConfig.URL, customName)
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

func (r *ShopRepositoryImpl) GetAllTags(shopID string) ([]entities.Tag, error) {
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
