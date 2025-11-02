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

type UserRepositoryImpl struct {
	db             database.Database
	supabaseConfig *config.Supabase
}

func NewUserRepositoryImpl(db database.Database, supabaseConfig *config.Supabase) UserRepository {
	return &UserRepositoryImpl{
		db:             db,
		supabaseConfig: supabaseConfig,
	}
}

func (r *UserRepositoryImpl) CreateUser(userModel *entities.User) error {
	return r.db.Connect().Create(userModel).Error
}

func (r *UserRepositoryImpl) CreateAdmin(adminModel *entities.ShopAdmin) error {
	return r.db.Connect().Create(adminModel).Error
}

func (r *UserRepositoryImpl) CreateShop(shopEntity *entities.Shop) error {
	return r.db.Connect().Create(shopEntity).Error
}

func (r *UserRepositoryImpl) CreateContact(contact *entities.Contact) error {
	return r.db.Connect().Create(contact).Error
}

func (r *UserRepositoryImpl) GetContactByID(contactID uuid.UUID) (*entities.Contact, error) {
	db := r.db.Connect()
	var contact entities.Contact
	if err := db.Where("contact_id = ?", contactID).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *UserRepositoryImpl) DeleteContact(contactID uuid.UUID) error {
	db := r.db.Connect()
	return db.Where("contact_id = ?", contactID).Delete(&entities.Contact{}).Error
}

func (r *UserRepositoryImpl) UploadUserProfileImage(username string, imageModel *utils.ImageModel) (string, error) {
	customName := username + "_" + fmt.Sprintf("%d", time.Now().Unix())

	mimeType := mime.TypeByExtension(imageModel.Ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	url := fmt.Sprintf("%s/storage/v1/object/images/userProfile/%s", r.supabaseConfig.URL, customName)
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
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/images/userProfile/%s", r.supabaseConfig.URL, customName)
	return publicURL, nil
}

func (r *UserRepositoryImpl) EditUserProfileImage(userID uuid.UUID, imageModel *utils.ImageModel) error {
	db := r.db.Connect()
	user, err := r.GetUserByUserID(userID)
	if err != nil {
		return err
	}

	if user.ProfileImageURL != "NULL" && user.ProfileImageURL != "" {

		publicPrefixRender := fmt.Sprintf("%s/storage/v1/render/image/public/", r.supabaseConfig.URL)
		objectPath := strings.TrimPrefix(user.ProfileImageURL, publicPrefixRender)

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

	newImageUrl, err := r.UploadUserProfileImage(user.Username, imageModel)
	if err != nil {
		return err
	}

	err = db.Model(&entities.User{}).
		Where("user_id = ?", userID).
		Update("profile_image_url", newImageUrl).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) GetUserByUsername(username string) (*entities.User, error) {
	var user entities.User
	db := r.db.Connect()
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByUserID(userID uuid.UUID) (*entities.User, error) {
	var user entities.User
	db := r.db.Connect()
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetShopByAdminID(adminID uuid.UUID) (*entities.Shop, error) {
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

func (r *UserRepositoryImpl) EditUser(userID uuid.UUID, user *entities.User) error {
	db := r.db.Connect()
	err := db.Model(&entities.User{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"f_name": user.FName,
			"l_name": user.LName,
			"gender": user.Gender,
		}).Error
	return err
}

func (r *UserRepositoryImpl) GetAllShops() ([]entities.Shop, error) {
	var shops []entities.Shop
	db := r.db.Connect()
	if err := db.Find(&shops).Error; err != nil {
		return nil, err
	}
	return shops, nil
}

func (r *UserRepositoryImpl) GetShopAdminByUsername(username string) (*entities.ShopAdmin, error) {
	var admin entities.ShopAdmin
	db := r.db.Connect()
	if err := db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *UserRepositoryImpl) Topup(topupModel *entities.TopUp) error {
	if err := r.db.Connect().Create(topupModel).Error; err != nil {
		return err
	}

	var user entities.User
	db := r.db.Connect()
	if err := db.First(&user, "user_id = ?", topupModel.UserID).Error; err != nil {
		return err
	}

	newBalance := user.Wallet + topupModel.Amount
	if err := db.Model(&entities.User{}).
		Where("user_id = ?", topupModel.UserID).
		Update("wallet", newBalance).Error; err != nil {
		return err
	}

	return nil
}
