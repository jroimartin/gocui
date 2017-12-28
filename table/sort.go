package table

type SortOrder int
type SortFn func(interface{}, interface{}) bool

const (
	SortNone SortOrder = iota
	SortAsc
	SortDesc
)

type SortBy struct {
	index  int
	order  SortOrder
	sortFn SortFn
}
