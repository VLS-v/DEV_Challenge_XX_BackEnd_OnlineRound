package saves

import (
	"encoding/json"
	"fmt"
	"os"
	_ "path/filepath"
	"spreadsheets/models"
	"time"
)

const SPREADSHEETS_FILE_NAME = "saves.json"

type Saves struct {
	SavesFile          *os.File
	SavesData          map[string]models.Sheet
	relativePath       string
	relativePathToFile string
}

func (sv *Saves) Open(filepath string) error {
	sv.relativePath = filepath
	sv.relativePathToFile = filepath + SPREADSHEETS_FILE_NAME

	file, err := os.OpenFile(sv.relativePathToFile, os.O_RDWR, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(sv.relativePathToFile)
			if err != nil {
				return fmt.Errorf("could not create file: %v", err)
			}
			fmt.Printf("File '%s' created.\n", sv.relativePathToFile)
			sv.SavesFile = file
			return nil
		}
		return fmt.Errorf("could not open file: %v", err)
	}
	fmt.Printf("File '%s' opened.\n", sv.relativePathToFile)
	sv.SavesFile = file
	return nil
}

func (sv *Saves) Load() error {
	fileInfo, err := sv.SavesFile.Stat()

	if err != nil {
		fmt.Println("Error retrieving file information.", err)
		return err
	}

	if fileInfo.Size() == 0 {
		sv.SavesData = make(map[string]models.Sheet)
		return nil
	}

	if err := json.NewDecoder(sv.SavesFile).Decode(&sv.SavesData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}

	return nil
}

func (sv *Saves) Write(newData map[string]models.Sheet) error {
	jsonData, err := json.Marshal(newData)
	if err != nil {
		return err
	}

	currentTime := time.Now().UnixNano()
	tempFileName := fmt.Sprintf("%s_tempSaves.%d.json", sv.relativePath, currentTime)
	sv.SavesFile.Close()

	err = os.Rename(sv.relativePathToFile, tempFileName)
	if err != nil {
		sv.Open(sv.relativePathToFile)
		sv.Load()
		return err
	}

	newSaves, err := os.Create(sv.relativePathToFile)
	if err != nil {
		os.Rename(tempFileName, sv.relativePathToFile)
		sv.Open(sv.relativePathToFile)
		sv.Load()
		return fmt.Errorf("could not create temp file: %v", err)
	}

	_, err = newSaves.Write(jsonData)
	if err != nil {
		newSaves.Close()
		os.Remove(sv.relativePathToFile)
		os.Rename(tempFileName, sv.relativePathToFile)
		return err
	}
	os.Remove(tempFileName)
	newSaves.Close()

	sv.Open(sv.relativePath)
	sv.Load()

	return nil
}
