package service

import (
	"bytes"
	"fmt"
	"hewhew-backend/config"
	"hewhew-backend/entities"
	"hewhew-backend/pkg/user/repository"
	"mime"
	"net/http"
	"time"
)

type UserServiceImpl struct {
	userRepository repository.UserRepository
	supabaseConfig *config.Supabase
}

func NewUserServiceImpl(userRepository repository.UserRepository, supabaseConfig *config.Supabase) UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
		supabaseConfig: supabaseConfig,
	}
}

func (s *UserServiceImpl) CreateUser(userEntity *entities.User) (*entities.User, error) {
	return s.userRepository.CreateUser(userEntity)
}

func (s *UserServiceImpl) GetUsers() ([]*entities.User, error) {
	return s.userRepository.GetUsers()
}

func (s *UserServiceImpl) UploadUserProfileImage(username string, ext string, image []byte) (string, error) {
	customName := username + "_" + fmt.Sprintf("%d", time.Now().Unix())

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	url := fmt.Sprintf("%s/storage/v1/object/images/userProfile/%s", s.supabaseConfig.URL, customName)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(image))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.supabaseConfig.Key))
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
	publicURL := fmt.Sprintf("%s/storage/v1/render/image/public/images/userProfile/%s", s.supabaseConfig.URL, customName)
	return publicURL, nil
}
