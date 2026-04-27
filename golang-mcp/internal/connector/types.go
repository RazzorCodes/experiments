package connector

type Content struct {
	Mime  string `json:"mime"`
	Value string `json:"value"`
}

type SearchResult struct {
	Title string    `json:"title"`
	Type  string    `json:"type"`
	Id    string    `json:"id"`
	Hints []Content `json:"hints"`
}

type ContentResult struct {
	Title   string  `json:"title"`
	Type    string  `json:"type"`
	Id      string  `json:"id"`
	Content Content `json:"content"`
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

type MoveResult struct {
	Title    string `json:"title"`
	Type     string `json:"type"`
	Id       string `json:"id"`
	BranchId string `json:"branch_id"`
}

type DeleteResult struct {
	Id string `json:"id"`
}
