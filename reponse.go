package pocketbase

type ResponseList struct {
	Page       int              `json:"page"`
	PerPage    int              `json:"perPage"`
	TotalItems int              `json:"totalItems"`
	TotalPages int              `json:"totalPages"`
	Items      []map[string]any `json:"items"`
}

type ResponseCreate struct {
	ID      string `json:"id"`
	Created string `json:"created"`
	Field   string `json:"field"`
	Updated string `json:"updated"`
}
