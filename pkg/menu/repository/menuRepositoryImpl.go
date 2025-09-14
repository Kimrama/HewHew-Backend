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
	"time"

	"github.com/google/uuid"
)
type MenuRepositoryImpl struct {
    db database.Database
    supabaseConfig *config.Supabase
}
func NewMenuRepositoryImpl(db database.Database, supabaseConfig *config.Supabase) MenuRepository {
return &MenuRepositoryImpl{
    db: db,
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
        body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload image: %s, %s", resp.Status, string(body))
	}
	publicURL := fmt.Sprintf("%s/storage/v1/render/image/public/images/userProfile/%s", r.supabaseConfig.URL, customName)
	return publicURL, nil
}   
