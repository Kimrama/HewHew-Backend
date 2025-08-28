package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func PreprocessUploadImage(file *multipart.FileHeader) ([]byte, string, error) {
	src, err := file.Open()
	if err != nil {
		return nil, "", err
	}
	defer src.Close()

	// decode image
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, "", err
	}

	// resize to 400x400
	resizedImg := resize.Resize(400, 400, img, resize.Lanczos3)

	// encode by filetype
	buf := new(bytes.Buffer)
	ext := strings.ToLower(filepath.Ext(file.Filename))

	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(buf, resizedImg, nil)
	case ".png":
		err = png.Encode(buf, resizedImg)
	default:
		err = fmt.Errorf("unsupported file type: %s", ext)
	}
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), ext, nil
}
