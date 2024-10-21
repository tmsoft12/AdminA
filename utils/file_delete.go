package utils

import (
	"fmt"
	"os"
	"time"
)

// DeleteFileWithRetry faýly pozmak üçin birnäçe synanyşyk edýär
func DeleteFileWithRetry(filePath string) error {
	// Pozmaga synanyşýarlar, 3 gezek gaýtalaýar
	for i := 0; i < 3; i++ {
		err := os.Remove(filePath)
		if err == nil {
			fmt.Println("Faýl üstünlikli pozuldy:", filePath)
			return nil
		}

		// Faýly pozup bolmasa, 2 sekunt garaşýar we gaýtadan synanyşýar
		fmt.Printf("Faýly pozmak synanyşygy %d başa barmady, ýene synanyşýar...\n", i+1)
		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("Faýly pozup bolmady: %s", filePath)
}
