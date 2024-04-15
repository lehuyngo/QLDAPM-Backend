package entities

type DraftContactReport struct {
	Total      int `json:"total"`
	Resolved   int `json:"resolved"`
	Processing int `json:"processing"`
	Deleted    int `json:"deleted"`
}
