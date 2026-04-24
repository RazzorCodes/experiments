package connectors

type SearchResult struct {
	Title string   `json:"title"`
	Type  string   `json:"type"`
	Id    string   `json:"id"`
	Hints []string `json:"hints"`
}
