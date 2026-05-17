package storage

import (
	"encoding/json"
	"os"
	"github.com/kaustubha-chaturvedi/Rook/internal/notebook"
)

func Save(path string, f notebook.File) error {
	b, err := json.MarshalIndent(f, "", "  ")

	if err != nil {
		return err
	}
	
	return os.WriteFile(path, b, 0o600)
}

func Load(path string) (notebook.File, error) {
	b, err := os.ReadFile(path)
	
	if err != nil {
		return notebook.File{}, err
	}
	
	var f notebook.File
	return f, json.Unmarshal(b, &f)
}
