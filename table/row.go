package table

type Row struct {
	table     *Table
	values    []interface{}
	strValues []string
}

type Rows []*Row

func (r Rows) Len() int {
	return len(r)
}

func (r Rows) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r Rows) Less(i, j int) bool {
	sortOrder := r[i].table.sort
	var k int
	for k = 0; k < len(sortOrder)-1; k++ {
		s := sortOrder[k]
		if s.order == SortDesc {
			if gt(r[i].values[s.index], r[j].values[s.index], s.sortFn) {
				return true
			}
			if gt(r[j].values[s.index], r[i].values[s.index], s.sortFn) {
				return false
			}
		} else {
			if lt(r[i].values[s.index], r[j].values[s.index], s.sortFn) {
				return true
			}
			if lt(r[j].values[s.index], r[i].values[s.index], s.sortFn) {
				return false
			}
		}
	}

	s := sortOrder[k]
	if s.order == SortDesc {
		return gt(r[i].values[s.index], r[j].values[s.index], s.sortFn)
	}
	return lt(r[i].values[s.index], r[j].values[s.index], s.sortFn)
}

func gt(a interface{}, b interface{}, fn SortFn) bool {
	return lt(b, a, fn)
}

func lt(a interface{}, b interface{}, fn SortFn) bool {
	if fn != nil {
		return fn(a, b)
	} else {
		switch a.(type) {
		case int:
			return a.(int) < b.(int)
		case string:
			return a.(string) < b.(string)
		}
	}
	return false
}
