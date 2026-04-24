package connectors

type SearchResult struct {
	Title string   `json:"title"`
	Type  string   `json:"type"`
	Id    string   `json:"id"`
	Hints []string `json:"hints"`
}

type ContentResult struct {
	Title   string `json:"title"`
	Type    string `json:"type"`
	Id      string `json:"id"`
	Content string `json:"content"`
}

type UpdateResult struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	Id    string `json:"id"`
}

type CreateResult struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}
