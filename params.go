package pocketbase

type ParamsList struct {
	Page    int
	Size    int
	Filters string
	Sort    string
	Expand  string

	hackResponseRef any //hack for collection list
}
