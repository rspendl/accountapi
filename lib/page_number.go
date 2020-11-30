package lib

// PageNumber is a convenience type for selecting the page number in the account List request.
// Beside regular unsigned numbers, it can have special values First and Last
// as the server API can use "first" and "last" keywords instead of numeric page numbers.
// To avoid using pointers for page numbers when the defaults are used, a special
// constant PNNone can be used to omit page number from List requests.
type PageNumber int

const (
	// First is used to fetch the "first" page instead of a numbered page.
	First PageNumber = -1
	// Last is used to fetch the "last" page instead of a numbered page.
	Last PageNumber = -2
	// PNNone signals that page number should be omitted from the List request.
	PNNone = -3
)
