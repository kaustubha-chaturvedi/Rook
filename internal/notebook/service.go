package notebook

import (
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type View struct {
	Name string `json:"name"`
	SQL  string `json:"sql"`
	Path string `json:"path,omitempty"`
}

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) ImportPath(path string) (View, error) {
	f, err := Load(path)

	if err != nil {
		return View{}, err
	}

	return View{
		Name: filepath.Base(path),
		SQL:  SQLFromFile(f),
		Path: path,
	}, nil
}

func (s *Service) ImportContent(name string, content string) (View, error) {
	f, err := ParseJSON([]byte(content))

	if err != nil {
		return View{}, err
	}

	return View{Name: name, SQL: SQLFromFile(f)}, nil
}

func (s *Service) PickAndImport() (View, error) {
	path, err := application.Get().Dialog.OpenFile().
		SetTitle("Open Notebook").
		AddFilter("SQL Notebook", "*.snb").
		PromptForSingleSelection()

	if err != nil {
		return View{}, err
	}

	if path == "" {
		return View{}, nil
	}

	return s.ImportPath(path)
}

func (s *Service) Save(path string, sql string) error {
	return Save(path, FileFromSQL(sql))
}

func (s *Service) PickAndSave(sql string) (string, error) {
	path, err := application.Get().Dialog.SaveFile().
		SetMessage("Save Notebook").
		AddFilter("SQL Notebook", "*.snb").
		SetFilename("notebook.snb").
		PromptForSingleSelection()

	if err != nil || path == "" {
		return "", err
	}
	
	return path, s.Save(path, sql)
}
