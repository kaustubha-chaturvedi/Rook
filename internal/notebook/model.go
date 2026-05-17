package notebook

type File struct {
	Version int    `json:"version"`
	Cells   []Cell `json:"cells"`
}

type Cell struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Connection string `json:"connection"`
	Query      string `json:"query"`
}
