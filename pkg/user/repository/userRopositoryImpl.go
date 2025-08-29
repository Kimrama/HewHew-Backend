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
	"time"
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

func (r *UserRepositoryImpl) GetUsers() ([]*entities.User, error) {
	var users []*entities.User
	r.db.Connect()
	return users, nil
}
