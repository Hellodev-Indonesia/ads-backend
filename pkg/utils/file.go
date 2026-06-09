package utils

import (
	"errors"
	"image"
	"image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

// ProcessBrandPhoto processes an uploaded image.
// It checks if size <= 2MB, resizes it to max 500x500 maintaining aspect ratio,
// saves it to "public/uploads/brands/" directory, and returns the public path.
func ProcessBrandPhoto(fileHeader *multipart.FileHeader, slug string) (string, error) {
	if fileHeader == nil {
		return "", errors.New("file is nil")
	}

	// 1. Check max size (2MB)
	if fileHeader.Size > 2*1024*1024 {
		return "", errors.New("file size exceeds 2MB limit")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 2. Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", errors.New("invalid image file")
	}

	// 3. Auto resize (max 500x500 maintaining aspect ratio)
	resizedImg := resize.Thumbnail(500, 500, img, resize.Lanczos3)

	// 4. Save to local storage
	// Ensure directory exists
	uploadDir := filepath.Join("public", "uploads", "brands")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", err
	}

	filename := slug + ".jpg" // We'll save it as jpg
	filePath := filepath.Join(uploadDir, filename)

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if err := jpeg.Encode(out, resizedImg, &jpeg.Options{Quality: 85}); err != nil {
		return "", err
	}

	// Return the public URL path
	return "/uploads/brands/" + filename, nil
}
