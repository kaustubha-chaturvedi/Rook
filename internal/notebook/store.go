package notebook

import (
	"encoding/json"
	"os"
)

func Load(path string) (File, error) {
	b, err := os.ReadFile(path)

	if err != nil {
		return File{}, err
	}

	return ParseJSON(b)
}

func Save(path string, f File) error {
	if f.Version == 0 {
		f.Version = CurrentVersion
	}

	b, err := json.MarshalIndent(f, "", "  ")
	
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0o600)
}

func ParseJSON(data []byte) (File, error) {
	var f File
	
	if err := json.Unmarshal(data, &f); err != nil {
		return File{}, err
	}

	if f.Version == 0 {
		f.Version = CurrentVersion
	}
	
	return f, nil
}
