package saves

import (
	"encoding/json"
	"fmt"
	"os"
	"spreadsheets/models"
)

type SavesInterface interface {
	//Create() error
	OpenSaves(filename string) error
	Read() error
	Write() error
	//Close() error
}

type Saves struct {
	SavesFile *os.File
	SavesData models.SavesData
}

type File struct {
	OsFile *os.File
}

func (sv *Saves) OpenSaves(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(filename)
			if err != nil {
				return nil, fmt.Errorf("could not create file: %v", err)
			}
			fmt.Printf("File '%s' created.\n", filename)
			return file, nil
		}
		return nil, fmt.Errorf("could not open file: %v", err)
	}
	fmt.Printf("File '%s' opened.\n", filename)
	return file, nil
}

func (sv *Saves) LoadSaves(savesFile *os.File) (models.SavesData, error) {
	var savesData models.SavesData = make(models.SavesData)
	fileInfo, err := savesFile.Stat()

	if err != nil {
		fmt.Println("Error retrieving file information.", err)
		return nil, err
	}

	if fileInfo.Size() == 0 {
		return savesData, nil
	}

	if err := json.NewDecoder(savesFile).Decode(&savesData); err != nil {
		fmt.Println("Error decoding JSON:0", err)
		return nil, err
	}

	return savesData, nil
}
