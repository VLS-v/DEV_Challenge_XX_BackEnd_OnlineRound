package saves

import (
	"encoding/json"
	"fmt"
	"os"
	"spreadsheets/models"
)

type SavesInterface interface {
	Open(filename string) error
	Load() error
	Write() error
}

type Saves struct {
	SavesFile *os.File
	SavesData map[string]models.Sheet
}

type File struct {
	OsFile *os.File
}

func (sv *Saves) Open(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filename)
			if err != nil {
				return fmt.Errorf("could not create file: %v", err)
			}
			fmt.Printf("File '%s' created.\n", filename)
			sv.SavesFile = file
			return nil
		}
		return fmt.Errorf("could not open file: %v", err)
	}
	fmt.Printf("File '%s' opened.\n", filename)
	sv.SavesFile = file
	return nil
}

func (sv *Saves) Load() error {
	//var savesData models.SavesData = make(models.SavesData)
	fileInfo, err := sv.SavesFile.Stat()

	if err != nil {
		fmt.Println("Error retrieving file information.", err)
		return err
	}

	if fileInfo.Size() == 0 {
		return nil
	}

	if err := json.NewDecoder(sv.SavesFile).Decode(&sv.SavesData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return err
	}
	return nil
}

func (sv *Saves) Write() error {
	jsonData, err := json.Marshal(sv.SavesData)
	if err != nil {
		return err
	}
	sv.SavesFile.Truncate(0)
	sv.SavesFile.Seek(0, 0)
	_, err = sv.SavesFile.Write(jsonData)

	if err != nil {
		return err
	}
	return nil
}
