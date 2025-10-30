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

type DropOffRepositoryImpl struct {
	db             database.Database
	supabaseConfig *config.Supabase
}

func NewDropOffRepositoryImpl(db database.Database, supabaseConfig *config.Supabase) DropOffRepository {
	return &DropOffRepositoryImpl{
		db:             db,
		supabaseConfig: supabaseConfig,
	}
}

func (dr *DropOffRepositoryImpl) CreateDropOff(dropOff *entities.DropOffLocation) error {
	return dr.db.Connect().Create(dropOff).Error
}

func (dr *DropOffRepositoryImpl) UploadDropOffImage(DropOffLocationID uuid.UUID, imageModel *utils.ImageModel) (string, error) {
	customName := DropOffLocationID.String() + "_" + fmt.Sprintf("%d", time.Now().Unix())

	mimeType := mime.TypeByExtension(imageModel.Ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	url := fmt.Sprintf("%s/storage/v1/object/images/dropoffImage/%s", dr.supabaseConfig.URL, customName)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(imageModel.Body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", dr.supabaseConfig.Key))
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
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/images/dropoffImage/%s", dr.supabaseConfig.URL, customName)
	return publicURL, nil
}

func (dr *DropOffRepositoryImpl) GetAllDropOffs() ([]*entities.DropOffLocation, error) {
	var dropOffs []*entities.DropOffLocation
	err := dr.db.Connect().Find(&dropOffs).Error
	if err != nil {
		return nil, err
	}
	return dropOffs, nil
}

func (dr *DropOffRepositoryImpl) GetDropOffByID(dropOffID uuid.UUID) (*entities.DropOffLocation, error) {
	var dropOff entities.DropOffLocation
	err := dr.db.Connect().Where("drop_off_id = ?", dropOffID).First(&dropOff).Error
	if err != nil {
		return nil, err
	}
	return &dropOff, nil
}
