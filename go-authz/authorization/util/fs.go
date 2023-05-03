package util

import (
	"authorization/config"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func SaveFileToLocal(fileName string, file *multipart.FileHeader) error {
	// create a new file for the avatar
	filePath := filepath.Join(config.StorageConfig.StaticRoot, config.StorageConfig.StaticAvatarPath, fileName)
	avatarFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer avatarFile.Close()

	// open the uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	// copy the contents of the uploaded file to the new file
	if _, err := io.Copy(avatarFile, uploadedFile); err != nil {
		return err
	}

	return nil
}

func DeleteFileInLocal(path string) error {
	if _, err := os.Stat(path); err == nil {
		// delete the file
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}
