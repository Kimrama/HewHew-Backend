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
	publicURL := fmt.Sprintf("%s/storage/v1/render/image/public/images/userProfile/%s", r.supabaseConfig.URL, customName)
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
