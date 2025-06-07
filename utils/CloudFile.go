package utils

import (
	"context"
	"fmt"
	"image"
	"os"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/disintegration/imaging"

	"github.com/joho/godotenv"
)

type Credentials struct {
	Name       string
	API_KEY    string
	API_SECRET string
}

func ReadCredentials() (Credentials, error) {
	var Data Credentials
	if err := godotenv.Load(); err != nil {
		return Data, err
	}
	Data.Name = os.Getenv("CLOUDINARY_NAME")
	Data.API_KEY = os.Getenv("CLOUDINARY_API_KEY")
	Data.API_SECRET = os.Getenv("CLOUDINARY_API_SECRET")
	return Data, nil
}

func ConnectCloudinary(data Credentials) *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromParams(data.Name, data.API_KEY, data.API_SECRET)
	if err != nil {
		fmt.Errorf("Cloudinary config error: %v", err)
	}
	return cld
}

func ReadImage(filepath string) (image.Image, error) {
	img, err := imaging.Open(filepath)
	return img, err
}

func CompressImageAndSaveit(img image.Image, path string) (string, error) {
	compressedPath := filepath.Join(os.TempDir(), "compressed_"+filepath.Base(path))
	err := imaging.Save(img, compressedPath, imaging.JPEGQuality(70))
	if err != nil {
		return "", fmt.Errorf("failed to save compressed image: %w", err)
	}
	return compressedPath, nil

}

func UploadImageCode(cld *cloudinary.Cloudinary, filePath string) (string, error) {
	resp, err := cld.Upload.Upload(
		context.Background(),
		filePath,
		uploader.UploadParams{Folder: "my_gallery"},
	)
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

func UploadAnImage(filePath string) (string, error) {
	var link string
	Data, err := ReadCredentials() //Data is credentials
	if err != nil {
		return link, err
	}
	img, err := ReadImage(filePath)
	if err != nil {
		return "", err
	}
	compressed_filePath, err := CompressImageAndSaveit(img, filePath)
	if err != nil {
		return "", err
	}

	cld := ConnectCloudinary(Data)
	url, err := UploadImageCode(cld, compressed_filePath)
	if err != nil {
		return "", err
	} else {
		return url, nil
	}

}
