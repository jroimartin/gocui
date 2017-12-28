package table

type SortOrder int

const (
	SortNone SortOrder = iota
	SortAsc
	SortDesc
)

type SortBy struct {
	index int
	order SortOrder
}
