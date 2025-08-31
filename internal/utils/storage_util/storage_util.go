package storage_util

import (
	"os"
)

func SaveFileInStorage(encryptedFile []byte, fldPathToSave string) error {
	if err := os.WriteFile(fldPathToSave, encryptedFile, 0644); err != nil {
		return err
	}

	return nil
}

func DelFileFromStorage(pathToDelFile string) error {
	if err := os.Remove(pathToDelFile); err != nil {
		return err
	}

	return nil
}
