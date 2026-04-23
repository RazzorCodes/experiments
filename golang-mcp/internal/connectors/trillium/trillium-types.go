package connectors

type NoteResult struct {
	Title    string `json:"title"`
	Type     string `json:"type"`
	Id       string `json:"id"`
	Contents string `json:"contents"`
}
