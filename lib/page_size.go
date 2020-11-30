package lib

// PageSize is a convenience type for selecting the page size in the account List request.
// To avoid using pointers for page size when the defaults are used, a special
// constant PSNone can be used to omit page size from List requests.
type PageSize int

const (
	// PSNone signals that page size should be omitted from the List request.
	PSNone = -1
)
